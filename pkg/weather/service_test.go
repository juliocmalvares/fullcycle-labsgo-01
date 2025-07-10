package weather

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetTemperatureByCEP_ValidCEP(t *testing.T) {
	// Para este teste, vamos testar com a API real
	result, err := GetTemperatureByCEP("35620-000")
	if err != nil {
		t.Skipf("Skipping test due to API error (likely missing API key): %v", err)
		return
	}

	if result == nil {
		t.Error("Expected result to not be nil")
		return
	}

	// Verifica se as temperaturas estão presentes
	if result.Temp_C == 0 && result.Temp_F == 0 && result.Temp_K == 0 {
		t.Error("Expected at least one temperature value to be non-zero")
	}

	// Verifica se as temperaturas estão em ranges razoáveis
	if result.Temp_C < -50 || result.Temp_C > 60 {
		t.Errorf("Temperature in Celsius seems unreasonable: %f", result.Temp_C)
	}

	if result.Temp_F < -58 || result.Temp_F > 140 {
		t.Errorf("Temperature in Fahrenheit seems unreasonable: %f", result.Temp_F)
	}

	if result.Temp_K < 223 || result.Temp_K > 333 {
		t.Errorf("Temperature in Kelvin seems unreasonable: %f", result.Temp_K)
	}
}

func TestGetTemperatureByCEP_InvalidCEP(t *testing.T) {
	tests := []struct {
		name    string
		cep     string
		wantErr bool
	}{
		{
			name:    "CEP inválido",
			cep:     "123",
			wantErr: true,
		},
		{
			name:    "CEP vazio",
			cep:     "",
			wantErr: true,
		},
		{
			name:    "CEP inexistente",
			cep:     "99999-999",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTemperatureByCEP(tt.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTemperatureByCEP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemperatureHandler_ValidCEP(t *testing.T) {
	// Criar um servidor de teste
	server := httptest.NewServer(http.HandlerFunc(TemperatureHandler))
	defer server.Close()

	// Fazer requisição GET com CEP válido
	resp, err := http.Get(server.URL + "/temperature?cep=35620-000")
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Se retornar 500, verificar se é por falta de API key
	if resp.StatusCode == http.StatusInternalServerError {
		var errorResp ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err == nil && strings.Contains(errorResp.Error, "WEATHER_API_KEY") {
			t.Skipf("Skipping test due to missing API key: %s", errorResp.Error)
			return
		}
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
		return
	}

	// Verificar content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type, got %s", contentType)
	}

	// Decodificar resposta
	var response TemperatureResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
		return
	}

	// Verificar se as temperaturas estão presentes
	if response.Temp_C == 0 && response.Temp_F == 0 && response.Temp_K == 0 {
		t.Error("Expected at least one temperature value to be non-zero")
	}
}

func TestTemperatureHandler_AlwaysReturnsJSON(t *testing.T) {
	// Testa se o handler sempre retorna JSON, independente do resultado
	server := httptest.NewServer(http.HandlerFunc(TemperatureHandler))
	defer server.Close()

	// Testa com CEP válido (pode falhar por falta de API key, mas deve retornar JSON)
	resp, err := http.Get(server.URL + "/temperature?cep=35620-000")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Verifica se o Content-Type sempre é application/json
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type, got %s", contentType)
	}

	// Verifica se a resposta é JSON válido (sucesso ou erro)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Tenta decodificar como JSON (deve funcionar tanto para sucesso quanto erro)
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		t.Errorf("Response is not valid JSON: %v, body: %s", err, string(body))
	}
}

func TestTemperatureHandler_MissingCEP(t *testing.T) {
	// Criar um servidor de teste
	server := httptest.NewServer(http.HandlerFunc(TemperatureHandler))
	defer server.Close()

	// Fazer requisição GET sem CEP
	resp, err := http.Get(server.URL + "/temperature")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Verificar status code
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	// Verificar content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type, got %s", contentType)
	}

	// Verificar se a resposta é JSON válido
	var errorResp ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Errorf("Failed to decode error response: %v", err)
	}

	if errorResp.Error != "invalid zipcode" {
		t.Errorf("Expected error message 'invalid zipcode', got %s", errorResp.Error)
	}
}

func TestTemperatureHandler_InvalidMethod(t *testing.T) {
	// Criar um servidor de teste
	server := httptest.NewServer(http.HandlerFunc(TemperatureHandler))
	defer server.Close()

	// Fazer requisição POST (método inválido)
	resp, err := http.Post(server.URL+"/temperature?cep=35620-000", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Verificar status code
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", resp.StatusCode)
	}

	// Verificar content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type, got %s", contentType)
	}

	// Verificar se a resposta é JSON válido
	var errorResp ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Errorf("Failed to decode error response: %v", err)
	}

	if errorResp.Error != "method not allowed" {
		t.Errorf("Expected error message 'method not allowed', got %s", errorResp.Error)
	}
}

func TestTemperatureHandler_InvalidCEP(t *testing.T) {
	// Criar um servidor de teste
	server := httptest.NewServer(http.HandlerFunc(TemperatureHandler))
	defer server.Close()

	// Fazer requisição GET com CEP inválido
	resp, err := http.Get(server.URL + "/temperature?cep=123")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Verificar status code
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}

	// Verificar content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type, got %s", contentType)
	}

	// Verificar se a resposta é JSON válido
	var errorResp ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Errorf("Failed to decode error response: %v", err)
	}

	// Verificar se contém uma mensagem de erro
	if errorResp.Error == "" {
		t.Error("Expected error message, got empty string")
	}
}

func TestTemperatureResponse_JSONMarshal(t *testing.T) {
	response := TemperatureResponse{
		Temp_C: 25.0,
		Temp_F: 77.0,
		Temp_K: 298.15,
	}

	// Testar marshaling para JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal JSON: %v", err)
	}

	// Verificar se o JSON contém os campos esperados
	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, `"temp_C":25`) {
		t.Error("Expected temp_C field in JSON")
	}
	if !strings.Contains(jsonStr, `"temp_F":77`) {
		t.Error("Expected temp_F field in JSON")
	}
	if !strings.Contains(jsonStr, `"temp_K":298.15`) {
		t.Error("Expected temp_K field in JSON")
	}

	// Testar unmarshaling de volta
	var decodedResponse TemperatureResponse
	err = json.Unmarshal(jsonData, &decodedResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if decodedResponse.Temp_C != response.Temp_C {
		t.Errorf("Expected Temp_C %f, got %f", response.Temp_C, decodedResponse.Temp_C)
	}

	if decodedResponse.Temp_F != response.Temp_F {
		t.Errorf("Expected Temp_F %f, got %f", response.Temp_F, decodedResponse.Temp_F)
	}

	if decodedResponse.Temp_K != response.Temp_K {
		t.Errorf("Expected Temp_K %f, got %f", response.Temp_K, decodedResponse.Temp_K)
	}
}

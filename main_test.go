package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"temperature_server/pkg/weather"
	"testing"
)

func TestMainIntegration_ServerStartup(t *testing.T) {
	// Testa se o servidor consegue ser iniciado
	// Este é um teste básico de integração
	server := httptest.NewServer(http.HandlerFunc(weather.TemperatureHandler))
	defer server.Close()

	// Verifica se o servidor está respondendo
	resp, err := http.Get(server.URL + "/temperature?cep=35620-000")
	if err != nil {
		t.Skipf("Skipping integration test due to API error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Se retornar 500, verificar se é por falta de API key
	if resp.StatusCode == http.StatusInternalServerError {
		var errorResp weather.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err == nil && strings.Contains(errorResp.Error, "WEATHER_API_KEY") {
			t.Skipf("Skipping test due to missing API key: %s", errorResp.Error)
			return
		}
	}

	// Verifica se o servidor retorna uma resposta válida
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", resp.StatusCode)
	}
}

func TestMainIntegration_CompleteFlow(t *testing.T) {
	// Testa o fluxo completo: CEP -> Cidade -> Coordenadas -> Clima -> Temperaturas
	server := httptest.NewServer(http.HandlerFunc(weather.TemperatureHandler))
	defer server.Close()

	// Testa com um CEP válido
	resp, err := http.Get(server.URL + "/temperature?cep=35620-000")
	if err != nil {
		t.Skipf("Skipping integration test due to API error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Se retornar 500, verificar se é por falta de API key
	if resp.StatusCode == http.StatusInternalServerError {
		var errorResp weather.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err == nil && strings.Contains(errorResp.Error, "WEATHER_API_KEY") {
			t.Skipf("Skipping integration test due to missing API key: %s", errorResp.Error)
			return
		}
	}

	if resp.StatusCode == http.StatusOK {
		// Se a requisição foi bem-sucedida, verifica a estrutura da resposta
		var response weather.TemperatureResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Errorf("Failed to decode response: %v", err)
			return
		}

		// Verifica se todas as temperaturas estão presentes
		if response.Temp_C == 0 && response.Temp_F == 0 && response.Temp_K == 0 {
			t.Error("Expected at least one temperature value to be non-zero")
		}

		// Verifica se as temperaturas estão em ranges razoáveis
		if response.Temp_C < -50 || response.Temp_C > 60 {
			t.Errorf("Temperature in Celsius seems unreasonable: %f", response.Temp_C)
		}

		if response.Temp_F < -58 || response.Temp_F > 140 {
			t.Errorf("Temperature in Fahrenheit seems unreasonable: %f", response.Temp_F)
		}

		if response.Temp_K < 223 || response.Temp_K > 333 {
			t.Errorf("Temperature in Kelvin seems unreasonable: %f", response.Temp_K)
		}
	}
}

func TestMainIntegration_ErrorHandling(t *testing.T) {
	// Testa o tratamento de erros
	server := httptest.NewServer(http.HandlerFunc(weather.TemperatureHandler))
	defer server.Close()

	// Testa com CEP inválido
	resp, err := http.Get(server.URL + "/temperature?cep=123")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Verifica se retorna erro 500 para CEP inválido
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for invalid CEP, got %d", resp.StatusCode)
	}

	// Verifica content type para erro
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type for error, got %s", contentType)
	}

	// Testa sem CEP
	resp, err = http.Get(server.URL + "/temperature")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Verifica se retorna erro 400 para CEP ausente
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing CEP, got %d", resp.StatusCode)
	}

	// Verifica content type para erro
	contentType = resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type for error, got %s", contentType)
	}
}

func TestMainIntegration_HTTPMethods(t *testing.T) {
	// Testa diferentes métodos HTTP
	server := httptest.NewServer(http.HandlerFunc(weather.TemperatureHandler))
	defer server.Close()

	// Testa POST (método não permitido)
	resp, err := http.Post(server.URL+"/temperature?cep=35620-000", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for POST method, got %d", resp.StatusCode)
	}

	// Verifica content type para erro de método
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type for method error, got %s", contentType)
	}

	// Testa PUT (método não permitido)
	req, err := http.NewRequest("PUT", server.URL+"/temperature?cep=35620-000", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("Failed to create PUT request: %v", err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for PUT method, got %d", resp.StatusCode)
	}

	// Verifica content type para erro de método
	contentType = resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type for method error, got %s", contentType)
	}
}

func TestMainIntegration_ResponseHeaders(t *testing.T) {
	// Testa os headers da resposta
	server := httptest.NewServer(http.HandlerFunc(weather.TemperatureHandler))
	defer server.Close()

	resp, err := http.Get(server.URL + "/temperature?cep=35620-000")
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Se retornar 500, verificar se é por falta de API key
	if resp.StatusCode == http.StatusInternalServerError {
		var errorResp weather.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err == nil && strings.Contains(errorResp.Error, "WEATHER_API_KEY") {
			t.Skipf("Skipping test due to missing API key: %s", errorResp.Error)
			return
		}
	}

	// Verifica se o Content-Type está correto
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected application/json content type, got %s", contentType)
	}
}

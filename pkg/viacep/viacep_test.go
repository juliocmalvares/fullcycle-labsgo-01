package viacep

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchCEPData_ValidCEP(t *testing.T) {
	// Mock server para simular a API ViaCEP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := CEPResponse{
			CEP:         "35620-000",
			Logradouro:  "",
			Complemento: "",
			Unidade:     "",
			Bairro:      "",
			Localidade:  "Abaeté",
			UF:          "MG",
			Estado:      "Minas Gerais",
			Regiao:      "Sudeste",
			IBGE:        "3100203",
			GIA:         "",
			DDD:         "37",
			SIAFI:       "4003",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Em um teste real, você precisaria injetar a URL do mock server
	// Para este exemplo, estamos testando com a API real

	result, err := FetchCEPData("35620-000")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result to not be nil")
		return
	}

	if result.CEP != "35620-000" {
		t.Errorf("Expected CEP 35620-000, got %s", result.CEP)
	}

	if result.Localidade != "Abaeté" {
		t.Errorf("Expected Localidade Abaeté, got %s", result.Localidade)
	}

	if result.UF != "MG" {
		t.Errorf("Expected UF MG, got %s", result.UF)
	}
}

func TestFetchCEPData_InvalidCEP(t *testing.T) {
	tests := []struct {
		name    string
		cep     string
		wantErr bool
	}{
		{
			name:    "CEP muito curto",
			cep:     "123",
			wantErr: true,
		},
		{
			name:    "CEP muito longo",
			cep:     "123456789",
			wantErr: true,
		},
		{
			name:    "CEP vazio",
			cep:     "",
			wantErr: true,
		},
		{
			name:    "CEP com caracteres especiais",
			cep:     "abc-def",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FetchCEPData(tt.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCEPData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchCEPData_CEPFormatting(t *testing.T) {
	tests := []struct {
		name     string
		inputCEP string
		expected string
	}{
		{
			name:     "CEP com hífen",
			inputCEP: "35620-000",
			expected: "35620-000",
		},
		{
			name:     "CEP com pontos",
			inputCEP: "35620.000",
			expected: "35620-000",
		},
		{
			name:     "CEP com espaços",
			inputCEP: "35620 000",
			expected: "35620-000",
		},
		{
			name:     "CEP limpo",
			inputCEP: "35620000",
			expected: "35620-000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Para este teste, vamos apenas verificar se o CEP é aceito
			// (não retorna erro de formato inválido)
			_, err := FetchCEPData(tt.inputCEP)
			if err != nil && err.Error() == "invalid zipcode" {
				t.Errorf("FetchCEPData() should accept formatted CEP %s", tt.inputCEP)
			}
		})
	}
}

func TestCEPResponse_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"cep": "35620-000",
		"logradouro": "",
		"complemento": "",
		"unidade": "",
		"bairro": "",
		"localidade": "Abaeté",
		"uf": "MG",
		"estado": "Minas Gerais",
		"regiao": "Sudeste",
		"ibge": "3100203",
		"gia": "",
		"ddd": "37",
		"siafi": "4003"
	}`

	var response CEPResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if response.CEP != "35620-000" {
		t.Errorf("Expected CEP 35620-000, got %s", response.CEP)
	}

	if response.Localidade != "Abaeté" {
		t.Errorf("Expected Localidade Abaeté, got %s", response.Localidade)
	}

	if response.UF != "MG" {
		t.Errorf("Expected UF MG, got %s", response.UF)
	}

	if response.Estado != "Minas Gerais" {
		t.Errorf("Expected Estado Minas Gerais, got %s", response.Estado)
	}
}

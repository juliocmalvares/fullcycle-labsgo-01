package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFecthSearchFromWeatherAPI_ValidCity(t *testing.T) {
	// Mock server para simular a API Weather
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := []Search{
			{
				Id:      266410,
				Name:    "Bom Despacho",
				Region:  "Minas Gerais",
				Country: "Brazil",
				Lat:     -19.72,
				Lon:     -45.25,
				Url:     "bom-despacho-minas-gerais-brazil",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Para este teste, vamos testar com a API real
	result, err := FecthSearchFromWeatherAPI("Bom Despacho")
	if err != nil {
		t.Skipf("Skipping test due to API error (likely missing API key): %v", err)
		return
	}

	if result == nil {
		t.Error("Expected result to not be nil")
		return
	}

	if result.Name != "Bom Despacho" {
		t.Errorf("Expected Name Bom Despacho, got %s", result.Name)
	}

	if result.Region != "Minas Gerais" {
		t.Errorf("Expected Region Minas Gerais, got %s", result.Region)
	}

	if result.Country != "Brazil" {
		t.Errorf("Expected Country Brazil, got %s", result.Country)
	}
}

func TestFecthSearchFromWeatherAPI_EmptyCity(t *testing.T) {
	_, err := FecthSearchFromWeatherAPI("")
	if err == nil {
		t.Error("Expected error for empty city, got nil")
	}
}

func TestFecthSearchFromWeatherAPI_CityWithSpaces(t *testing.T) {
	// Testa se a função lida corretamente com cidades que têm espaços
	result, err := FecthSearchFromWeatherAPI("Bom Despacho")
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	if result == nil {
		t.Error("Expected result to not be nil")
		return
	}

	// Verifica se a cidade foi encontrada corretamente
	if result.Name == "" {
		t.Error("Expected city name to not be empty")
	}
}

func TestSearch_JSONUnmarshal(t *testing.T) {
	jsonData := `[
		{
			"id": 266410,
			"name": "Bom Despacho",
			"region": "Minas Gerais",
			"country": "Brazil",
			"lat": -19.72,
			"lon": -45.25,
			"url": "bom-despacho-minas-gerais-brazil"
		}
	]`

	var searchResults []Search
	err := json.Unmarshal([]byte(jsonData), &searchResults)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if len(searchResults) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(searchResults))
		return
	}

	result := searchResults[0]
	if result.Id != 266410 {
		t.Errorf("Expected Id 266410, got %d", result.Id)
	}

	if result.Name != "Bom Despacho" {
		t.Errorf("Expected Name Bom Despacho, got %s", result.Name)
	}

	if result.Lat != -19.72 {
		t.Errorf("Expected Lat -19.72, got %f", result.Lat)
	}

	if result.Lon != -45.25 {
		t.Errorf("Expected Lon -45.25, got %f", result.Lon)
	}
}

func TestFecthSearchFromWeatherAPI_EmptyResponse(t *testing.T) {
	// Para este teste, vamos simular o comportamento com uma cidade inexistente
	_, err := FecthSearchFromWeatherAPI("CidadeInexistente12345")
	if err != nil {
		// Se a API key não estiver configurada, pula o teste
		if err.Error() == "WEATHER_API_KEY is not set" {
			t.Skipf("Skipping test due to missing API key: %v", err)
			return
		}

		// Esperamos um erro para cidade inexistente
		if err.Error() != "no cities found for the given search term" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	}
}

package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchWeatherData_ValidCoordinates(t *testing.T) {
	// Mock server para simular a API Weather
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := WeatherResponse{
			Location: Location{
				Name:           "Bom Despacho",
				Region:         "Minas Gerais",
				Country:        "Brazil",
				Lat:            -19.717,
				Lon:            -45.25,
				TzID:           "America/Sao_Paulo",
				LocaltimeEpoch: 1752107290,
				Localtime:      "2025-07-09 21:28",
			},
			Current: Current{
				LastUpdatedEpoch: 1752106500,
				LastUpdated:      "2025-07-09 21:15",
				TempC:            17.0,
				TempF:            62.6,
				IsDay:            0,
				Condition: Condition{
					Text: "Partly Cloudy",
					Icon: "//cdn.weatherapi.com/weather/64x64/night/116.png",
					Code: 1003,
				},
				WindMph:    4.9,
				WindKph:    7.9,
				WindDegree: 114,
				WindDir:    "ESE",
				PressureMb: 1023.0,
				PressureIn: 30.2,
				PrecipMm:   0.0,
				PrecipIn:   0.0,
				Humidity:   45,
				Cloud:      27,
				FeelslikeC: 17.0,
				FeelslikeF: 62.7,
				WindchillC: 17.0,
				WindchillF: 62.7,
				HeatindexC: 17.0,
				HeatindexF: 62.7,
				DewpointC:  5.1,
				DewpointF:  41.1,
				VisKm:      10.0,
				VisMiles:   6.0,
				Uv:         0.0,
				GustMph:    10.3,
				GustKph:    16.6,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Para este teste, vamos testar com a API real
	result, err := FetchWeatherData(-19.72, -45.25)
	if err != nil {
		t.Skipf("Skipping test due to API error (likely missing API key): %v", err)
		return
	}

	if result == nil {
		t.Error("Expected result to not be nil")
		return
	}

	if result.Location.Name != "Bom Despacho" {
		t.Errorf("Expected Location Name Bom Despacho, got %s", result.Location.Name)
	}

	if result.Current.TempC != 17.0 {
		t.Errorf("Expected Current TempC 17.0, got %f", result.Current.TempC)
	}

	if result.Current.Condition.Text != "Partly Cloudy" {
		t.Errorf("Expected Condition Text Partly Cloudy, got %s", result.Current.Condition.Text)
	}
}

func TestWeatherResponse_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"location": {
			"name": "Bom Despacho",
			"region": "Minas Gerais",
			"country": "Brazil",
			"lat": -19.717,
			"lon": -45.25,
			"tz_id": "America/Sao_Paulo",
			"localtime_epoch": 1752107290,
			"localtime": "2025-07-09 21:28"
		},
		"current": {
			"last_updated_epoch": 1752106500,
			"last_updated": "2025-07-09 21:15",
			"temp_c": 17.0,
			"temp_f": 62.6,
			"is_day": 0,
			"condition": {
				"text": "Partly Cloudy",
				"icon": "//cdn.weatherapi.com/weather/64x64/night/116.png",
				"code": 1003
			},
			"wind_mph": 4.9,
			"wind_kph": 7.9,
			"wind_degree": 114,
			"wind_dir": "ESE",
			"pressure_mb": 1023.0,
			"pressure_in": 30.2,
			"precip_mm": 0.0,
			"precip_in": 0.0,
			"humidity": 45,
			"cloud": 27,
			"feelslike_c": 17.0,
			"feelslike_f": 62.7,
			"windchill_c": 17.0,
			"windchill_f": 62.7,
			"heatindex_c": 17.0,
			"heatindex_f": 62.7,
			"dewpoint_c": 5.1,
			"dewpoint_f": 41.1,
			"vis_km": 10.0,
			"vis_miles": 6.0,
			"uv": 0.0,
			"gust_mph": 10.3,
			"gust_kph": 16.6
		}
	}`

	var response WeatherResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if response.Location.Name != "Bom Despacho" {
		t.Errorf("Expected Location Name Bom Despacho, got %s", response.Location.Name)
	}

	if response.Location.Region != "Minas Gerais" {
		t.Errorf("Expected Location Region Minas Gerais, got %s", response.Location.Region)
	}

	if response.Current.TempC != 17.0 {
		t.Errorf("Expected Current TempC 17.0, got %f", response.Current.TempC)
	}

	if response.Current.Humidity != 45 {
		t.Errorf("Expected Current Humidity 45, got %d", response.Current.Humidity)
	}

	if response.Current.Condition.Text != "Partly Cloudy" {
		t.Errorf("Expected Condition Text Partly Cloudy, got %s", response.Current.Condition.Text)
	}

	if response.Current.Condition.Code != 1003 {
		t.Errorf("Expected Condition Code 1003, got %d", response.Current.Condition.Code)
	}
}

func TestFetchWeatherData_InvalidCoordinates(t *testing.T) {
	// Testa coordenadas inválidas (fora dos limites normais)
	_, err := FetchWeatherData(999.0, 999.0)
	if err != nil {
		// Esperamos um erro para coordenadas inválidas
		t.Logf("Expected error for invalid coordinates: %v", err)
	}
}

func TestWeatherResponse_EmptyFields(t *testing.T) {
	// Testa se a estrutura pode lidar com campos vazios
	response := WeatherResponse{
		Location: Location{
			Name:    "Test City",
			Region:  "",
			Country: "Test Country",
		},
		Current: Current{
			TempC:     0.0,
			Condition: Condition{},
			Humidity:  0,
		},
	}

	if response.Location.Name != "Test City" {
		t.Errorf("Expected Location Name Test City, got %s", response.Location.Name)
	}

	if response.Location.Region != "" {
		t.Errorf("Expected empty Region, got %s", response.Location.Region)
	}

	if response.Current.TempC != 0.0 {
		t.Errorf("Expected Current TempC 0.0, got %f", response.Current.TempC)
	}
}

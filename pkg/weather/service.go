package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"temperature_server/pkg/utils"
	"temperature_server/pkg/viacep"
)

type TemperatureResponse struct {
	Temp_C float64 `json:"temp_C"`
	Temp_F float64 `json:"temp_F"`
	Temp_K float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func GetTemperatureByCEP(cep string) (*TemperatureResponse, error) {
	cepData, err := viacep.FetchCEPData(cep)
	if err != nil {
		return nil, fmt.Errorf("can not find zipcode: %w", err)
	}

	searchData, err := FecthSearchFromWeatherAPI(cepData.Localidade)
	if err != nil {
		return nil, fmt.Errorf("can not find city: %w", err)
	}

	weatherData, err := FetchWeatherData(searchData.Lat, searchData.Lon)
	if err != nil {
		return nil, fmt.Errorf("can not find zipcode: %w", err)
	}

	formattedTemp := utils.FormatTemperatures(weatherData.Current.TempC)

	response := &TemperatureResponse{
		Temp_C: formattedTemp["temp_C"],
		Temp_F: formattedTemp["temp_F"],
		Temp_K: formattedTemp["temp_K"],
	}

	return response, nil
}

func TemperatureHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	cep := r.URL.Query().Get("cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid zipcode"})
		return
	}

	response, err := GetTemperatureByCEP(cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

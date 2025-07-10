package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type Search struct {
	Id      uint    `json:"id"`
	Name    string  `json:"name"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Url     string  `json:"url"`
}

func FecthSearchFromWeatherAPI(city string) (*Search, error) {

	city = strings.ReplaceAll(city, " ", "+")
	url := fmt.Sprintf("https://api.weatherapi.com/v1/search.json?q=%s", city)
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	apiKey := viper.GetString("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("WEATHER_API_KEY is not set")
	}
	req.Header.Set("key", apiKey)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	fmt.Println(response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var search []Search
	err = json.Unmarshal(body, &search)
	if err != nil {
		return nil, err
	}

	if len(search) == 0 {
		return nil, errors.New("no cities found for the given search term")
	}

	return &search[0], nil
}

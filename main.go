package main

import (
	"fmt"
	"log"
	"net/http"
	"temperature_server/pkg/weather"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	http.HandleFunc("/temperature", weather.TemperatureHandler)

	port := ":8080"
	fmt.Printf("Servidor rodando na porta %s\n", port)
	fmt.Printf("Teste com: http://localhost%s/temperature?cep=35620-000\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}

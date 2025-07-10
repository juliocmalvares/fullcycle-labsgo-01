package viacep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
}

func FetchCEPData(cep string) (*CEPResponse, error) {
	cep = strings.ReplaceAll(cep, "-", "")
	cep = strings.ReplaceAll(cep, ".", "")
	cep = strings.ReplaceAll(cep, " ", "")

	if len(cep) != 8 {
		return nil, errors.New("invalid zipcode")
	}

	var result CEPResponse
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep)
	fmt.Println(url)
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return &result, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &result, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &result, err
	}
	err = json.Unmarshal(body, &result)
	fmt.Println(result)
	return &result, err
}

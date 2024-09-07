package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type APIResponse interface{}

type ApiResponse struct {
	Api string `json:"Api"`
}

type ViaCEP struct {
	APIResponse
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

type ApiCEP struct {
	APIResponse
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

var viacep string = "97061306"
var apicep string = "98130000"

func main() {
	channelViaCEP := make(chan ViaCEP)
	channelApiCEP := make(chan ApiCEP)

	go GetViaCEP(channelViaCEP)
	go GetApiCEP(channelApiCEP)

	select {
	case resViaCEP := <-channelViaCEP:
		fmt.Printf("ViaCEP: %+v\n", resViaCEP)

	case resApiCEP := <-channelApiCEP:
		fmt.Printf("ApiCEP: %+v\n", resApiCEP)

	case <-time.After(time.Second):
		fmt.Printf("TimeOut")

	}

}

func GetViaCEP(chApi chan ViaCEP) {
	var viaCEP ViaCEP
	BuscaCep("https://viacep.com.br/ws/"+viacep+"/json/", &viaCEP)
	viaCEP.APIResponse = "ViaCEP"
	chApi <- viaCEP
}

func GetApiCEP(chVia chan ApiCEP) {
	var apiCEP ApiCEP
	BuscaCep("https://cdn.apicep.com/file/apicep/"+apicep+".json", &apiCEP)
	apiCEP.APIResponse = "ApiCEP"
	chVia <- apiCEP
}

func BuscaCep(url string, res APIResponse) error {
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}
	return nil

}

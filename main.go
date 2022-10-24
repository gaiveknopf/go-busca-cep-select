package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	fmt.Println("Buscando endereço, aguarde...")
	c1 := make(chan ViaCEP)
	c2 := make(chan ApiCEP)

	for _, cep := range os.Args[1:] {
		go func() {
			c1 <- getAddressByViaCEP(cep)
		}()

		go func() {
			c2 <- getAddressByApiCEP(cep)
		}()

	}
	select {
	case viaCep := <-c1:
		printViaCep(viaCep)
	case apiCep := <-c2:
		printApiCep(apiCep)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}

func getViaCEPUrl(cep string) string {
	return fmt.Sprint("https://viacep.com.br/ws/", cep, "/json/")
}

func getApiCEPUrl(cep string) string {
	return fmt.Sprint("https://cdn.apicep.com/file/apicep/", cep, ".json")
}

func getAddressByViaCEP(cep string) ViaCEP {
	req, err := http.Get(getViaCEPUrl(cep))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		panic(err)
	}
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
		panic(err)
	}
	var data ViaCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
		panic(err)
	}
	req.Body.Close()
	return data
}

func getAddressByApiCEP(cep string) ApiCEP {
	if len(cep) == 8 {
		firstFiveDigits := cep[:5]
		lastThreeDigits := cep[len(cep)-3:]
		cep = fmt.Sprint(firstFiveDigits, "-", lastThreeDigits)
	}
	req, err := http.Get(getApiCEPUrl(cep))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		panic(err)
	}
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
		panic(err)
	}
	var data ApiCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
		panic(err)
	}
	req.Body.Close()
	return data
}

func printViaCep(data ViaCEP) {
	fmt.Printf(
		`
API: %s,
PATH: https://viacep.com.br/ws/%s/json/,
CEP: %s,
Logradouro: %s,
Complemento: %s,
Bairro: %s,
Localidade: %s,
UF: %s,
IBGE: %s,
GIA: %s,
DDD: %s,
SIAFI: %s`,
		"ViaCEP",
		data.Cep,
		data.Cep,
		data.Logradouro,
		data.Complemento,
		data.Bairro,
		data.Localidade,
		data.Uf,
		data.Ibge,
		data.Gia,
		data.Ddd,
		data.Siafi)
}

func printApiCep(data ApiCEP) {
	fmt.Printf(
		`
API: %s,
PATH: https://cdn.apicep.com/file/apicep/%s.json,
CEP: %s,
Logradouro: %s,
Bairro: %s,
Localidade: %s,
UF: %s,
Status: %d,
Ok: %v,
StatusText: %s`,
		"ApiCEP",
		data.Code,
		data.Code,
		data.Address,
		data.District,
		data.City,
		data.State,
		data.Status,
		data.Ok,
		data.StatusText)
}

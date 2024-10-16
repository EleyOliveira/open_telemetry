package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type CepRequest struct {
	Cep string `json:"cep"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/consulta", consultaCep)
	println("rodando na porta 8080")
	http.ListenAndServe(":8080", nil)

}

func consultaCep(w http.ResponseWriter, r *http.Request) {

	//cria um limite para evitar ler um valor de requisição muito grande
	bodyReader := io.LimitReader(r.Body, 1048576)

	var cep CepRequest
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	valido, err := ValidarCep(cep.Cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ocorreu um erro ao validar o CEP: %s\n", err.Error())
		return
	}

	if !valido {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, "invalid zipcode")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cep)
}

func ValidarCep(cep string) (bool, error) {

	valido, err := regexp.MatchString(`^\d{8}$`, cep)

	if err != nil {
		return false, err
	}

	return valido, nil
}

package web

import (
	"curso/labs/open_telemetry/CepService/internal/infra/entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/index.html")
}

func ConsultaCep(w http.ResponseWriter, r *http.Request) {

	//cria um limite para evitar ler um valor de requisição muito grande
	bodyReader := io.LimitReader(r.Body, 1048576)

	var cep entity.Cep
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

	cepResponse, err := ValidarCep(cep.Cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ocorreu um erro ao consultar o CEP: %s\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cepResponse)
}

func ValidarCep(cep string) (bool, error) {

	valido, err := regexp.MatchString(`^\d{8}$`, cep)

	if err != nil {
		return false, err
	}

	return valido, nil
}

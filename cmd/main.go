package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

	// Log the raw body for debugging
	log.Println("Raw Request Body:", string(body))

	err = json.Unmarshal(body, &cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cep)
}

package main

import (
	"curso/labs/open_telemetry/CepService/internal/infra/web"
	"net/http"
)

func main() {

	http.HandleFunc("/", web.Index)
	http.HandleFunc("/consulta", web.ConsultaCep)

	println("rodando na porta 8080")

	http.ListenAndServe(":8080", nil)
}

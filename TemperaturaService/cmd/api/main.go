package main

import (
	"fmt"
	"net/http"
	"opentelemetry_zipkin/TemperaturaService/internal/infra/web"
)

func main() {
	http.HandleFunc("/temperatura", web.ConsultaTemperatura)
	fmt.Println("Servidor iniciado na porta 8081")
	http.ListenAndServe(":8081", nil)
}

package main

import (
	"curso/labs/open_telemetry/TemperaturaService/internal/infra/web"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/temperatura", web.ConsultaTemperatura)
	fmt.Println("Servidor iniciado na porta 8081")
	http.ListenAndServe(":8081", nil)
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"opentelemetry_zipkin/TemperaturaService/internal/infra/web"
	"opentelemetry_zipkin/pkg"
)

func main() {

	// Inicialize o TracerProvider
	tp, err := pkg.InitTracer("TemperaturaService")
	if err != nil {
		fmt.Printf("Failed to initialize tracer: %v\n", err)
		return
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}()

	http.HandleFunc("/temperatura", web.ConsultaTemperatura)
	fmt.Println("Servidor iniciado na porta 8081")
	http.ListenAndServe(":8081", nil)
}

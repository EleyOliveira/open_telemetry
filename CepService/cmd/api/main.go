package main

import (
	"context"
	"fmt"
	"net/http"
	"opentelemetry_zipkin/CepService/internal/infra/web"
	"opentelemetry_zipkin/pkg"
)

func main() {

	// Inicialize o TracerProvider
	tp, err := pkg.InitTracer("CepService")
	if err != nil {
		fmt.Printf("Failed to initialize tracer: %v\n", err)
		return
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}()

	http.HandleFunc("/", web.Index)
	http.HandleFunc("/consulta", web.ConsultaTemperatura)

	println("rodando na porta 8080")

	http.ListenAndServe(":8080", nil)
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"opentelemetry_zipkin/CepService/internal/infra/web"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {

	// Inicialize o TracerProvider
	tp, err := initTracer()
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

func initTracer() (*trace.TracerProvider, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// Configure atributos do serviço, como nome, versão, etc.
			attribute.String("service.name", "CepService"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Configure o exportador OTLP para enviar dados para o coletor.
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Crie o TracerProvider.
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

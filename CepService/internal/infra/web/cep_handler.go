package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"opentelemetry_zipkin/CepService/internal/infra/entity"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func Index(w http.ResponseWriter, r *http.Request) {

	// Obtem o tracer global configurado no arquivo main.go
	tr := otel.Tracer("CepService")

	//Inicia um novo span. A chamada para EndSpan deve ser a última ação da função
	_, span := tr.Start(r.Context(), "Index")
	defer span.End()

	http.ServeFile(w, r, "../static/index.html")
}

func ConsultaTemperatura(w http.ResponseWriter, r *http.Request) {

	tr := otel.Tracer("CepService")
	ctx, span := tr.Start(r.Context(), "ConsultaTemperatura")
	defer span.End()

	span.SetAttributes(attribute.String(
		"http.method",
		r.Method,
	))

	span.SetAttributes(attribute.String(
		"http.url",
		r.URL.String(),
	))

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

	//resp, err := http.Get(fmt.Sprintf("https://temperatura-l5x7giwwma-uc.a.run.app/cep?cep=%s", cep.Cep))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://temperatura-l5x7giwwma-uc.a.run.app/cep?cep=%s", cep.Cep), nil)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ocorreu um erro ao preparar a requisição: %s\n", err.Error())
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ocorreu um erro ao consultar a temperatura: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ocorreu um erro ao ler a resposta: %s\n", err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintf(w, "Erro ao consultar a temperatura: %s\n", string(body))
		return
	}

	var temperatura entity.Temperature
	err = json.Unmarshal(body, &temperatura)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ocorreu um erro ao decodificar a temperatura: %s\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temperatura)
}

func ValidarCep(cep string) (bool, error) {

	valido, err := regexp.MatchString(`^\d{8}$`, cep)

	if err != nil {
		return false, err
	}

	return valido, nil
}

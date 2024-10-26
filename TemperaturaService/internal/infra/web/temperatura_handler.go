package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ViaCEP struct {
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

type WeatherTemperature struct {
	Location struct {
		Name string `json:"name"`
	}
	Current struct {
		TempC float64 `json:"temp_c"`
	}
}

type Temperature struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF string  `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func ConsultaTemperatura(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	carrier := propagation.HeaderCarrier(r.Header)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	tr := otel.Tracer("TemperaturaService")
	ctx, span := tr.Start(ctx, "ConsultaTemperatura")
	defer span.End()

	cep := r.URL.Query().Get("cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "CEP não informado")
		return
	}

	data, codigo, err := GetCep(ctx, cep)
	if err != nil {
		w.WriteHeader(codigo)
		fmt.Fprintf(w, "ocorreu o erro: %s", err.Error())
		return
	}

	dataTemperature, codigo, err := GetTemperature(ctx, data.Localidade)
	if err != nil {
		w.WriteHeader(codigo)
		fmt.Fprintf(w, "ocorreu o erro: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataTemperature)
}

func GetCep(ctx context.Context, cep string) (*ViaCEP, int, error) {

	tr := otel.Tracer("TemperaturaService")
	_, span := tr.Start(ctx, "GetCep")
	defer span.End()

	req, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erro ao fazer requisição da api de CEP: %s", err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erro ao ler resposta do CEP: %s", err)
	}

	var data ViaCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erro ao formatar a resposta: %s", err)
	}

	if data.Erro == "true" {
		return nil, http.StatusNotFound, errors.New("can not find zipcode")
	}

	return &data, http.StatusOK, nil
}

func GetTemperature(ctx context.Context, localidade string) (*Temperature, int, error) {

	tr := otel.Tracer("TemperaturaService")
	_, span := tr.Start(ctx, "GetTemperature")
	defer span.End()

	apiKey := "90afc375b7bf4a7cb18171824242909"
	params := url.Values{}
	params.Add("q", localidade)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&%s", apiKey, params.Encode())

	req, err := http.Get(url)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erro ao fazer requisição da api de temperatura: %s", err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erro ao ler resposta da temperatura: %s", err)
	}

	var dataWeather WeatherTemperature
	err = json.Unmarshal(res, &dataWeather)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erro ao formatar a resposta da temperatura: %s", err)
	}

	dataTemperature := Temperature{}
	dataTemperature.City = dataWeather.Location.Name
	dataTemperature.TempC = dataWeather.Current.TempC
	dataTemperature.ConverteCelsiusFarenheit(dataWeather.Current.TempC)
	dataTemperature.ConverteCelsiusKelvin(dataWeather.Current.TempC)

	return &dataTemperature, http.StatusOK, nil

}

func (t *Temperature) ConverteCelsiusFarenheit(grauCelsius float64) {
	t.TempF = fmt.Sprintf("%.1f", grauCelsius*1.8+32)
}

func (t *Temperature) ConverteCelsiusKelvin(grauCelsius float64) {
	t.TempK = grauCelsius + 273
}

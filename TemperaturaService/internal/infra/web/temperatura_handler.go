package web

import (
	"curso/labs/open_telemetry/TemperaturaService/internal/infra/entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ConsultaTemperatura(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "CEP não informado")
		return
	}

	fmt.Println("CEP recebido:", cep)

	// Simulando a consulta a um banco de dados ou API externa
	temperatura, err := obterTemperatura(cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Erro ao consultar a temperatura: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temperatura)
}

func obterTemperatura(cep string) (*entity.Temperature, error) {
	resp, err := http.Get(fmt.Sprintf("https://temperatura-l5x7giwwma-uc.a.run.app/cep?cep=%s", cep))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var temperatura entity.Temperature
	err = json.Unmarshal(body, &temperatura)
	if err != nil {
		return nil, err
	}

	return &temperatura, nil

}

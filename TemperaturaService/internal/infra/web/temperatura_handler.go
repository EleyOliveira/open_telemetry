package web

import (
	"fmt"
	"net/http"
	"strconv"
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
	fmt.Fprintf(w, "Temperatura: %d°C", temperatura)
}

func obterTemperatura(cep string) (int, error) {
	// Lógica para consultar a temperatura com base no CEP
	// Neste exemplo, estamos simulando a temperatura com base no CEP
	// Em um cenário real, você consultaria um banco de dados ou API externa
	temperatura, err := strconv.Atoi(cep)
	if err != nil {
		return 0, err
	}
	return temperatura, nil
}

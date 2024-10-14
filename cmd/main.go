package main

import "net/http"

type CepRequest struct {
	Cep string `json:"cep"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/consulta", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("cep legal"))
	})

	http.ListenAndServe(":8080", nil)
	println("rodando na porta 8080")
}

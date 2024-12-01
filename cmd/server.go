package main

import (
	"context"
	"encoding/json"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/adapter"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/application"
	"net/http"
)

type ResultError struct {
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{cep}", func(w http.ResponseWriter, r *http.Request) {
		cep := r.PathValue("cep")
		coordinateRepository := adapter.NewCoordinateRepository()
		weatherRepository := adapter.NewWeatherRepository()
		w.Header().Set("Content-Type", "application/json")
		s := application.NewWeatherService(coordinateRepository, weatherRepository)
		output, e := s.GetTemperature(context.Background(), cep)

		if e != nil {
			switch e {
			case application.ErrCepInvalid, application.ErrCepMalformed:
				w.WriteHeader(http.StatusUnprocessableEntity)
				_ = json.NewEncoder(w).Encode(&ResultError{Message: "invalid zipcode"})
			case application.ErrCepNotFound:
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(&ResultError{Message: "can not find zipcode"})
			default:
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(&ResultError{Message: "internal server error" + e.Error()})
			}
			return
		}

		_ = json.NewEncoder(w).Encode(output)
	})
	_ = http.ListenAndServe(":8080", mux)
}

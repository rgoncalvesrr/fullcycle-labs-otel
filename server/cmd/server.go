package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/adapter"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/application"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/configs"
	"net/http"
	"time"
)

type ResultError struct {
	Message string `json:"message"`
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Handle("/metrics", promhttp.Handler())

	router.Get("/{cep}", handler)

	addr := fmt.Sprintf(":%s", configs.Cfg.ServerApiPort)

	_ = http.ListenAndServe(addr, router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
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
}

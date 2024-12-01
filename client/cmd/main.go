package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/client/configs"
	"net"
	"net/http"
	"regexp"
	"time"
)

type WeatherOutput struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_k"`
}

type Input struct {
	Cep string `json:"cep"`
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Handle("/metrics", promhttp.Handler())

	router.Post("/", handler)

	addr := fmt.Sprintf(":%s", configs.Cfg.ClientApiHttpPort)

	_ = http.ListenAndServe(addr, router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var i Input
	e := json.NewDecoder(r.Body).Decode(&i)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	match, _ := regexp.MatchString("[0-9]{8}", i.Cep)

	w.Header().Set("Content-Type", "application/json")
	if !match {
		msg := "invalid zipcode"
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}

	uri := fmt.Sprintf("http://%s/%s", net.JoinHostPort(configs.Cfg.ServerApiHost, configs.Cfg.ServerApiPort), i.Cep)
	client := resty.New()
	req, e := client.R().
		SetContext(r.Context()).
		SetHeader("Accept", "application/json").
		SetResult(&WeatherOutput{}).
		Get(uri)

	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	result := &WeatherOutput{
		City:       req.Result().(*WeatherOutput).City,
		Celsius:    req.Result().(*WeatherOutput).Celsius,
		Fahrenheit: req.Result().(*WeatherOutput).Fahrenheit,
		Kelvin:     req.Result().(*WeatherOutput).Kelvin,
	}

	json.NewEncoder(w).Encode(result)
}

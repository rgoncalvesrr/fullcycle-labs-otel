package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/configs"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/pkg/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net"
	"net/http"
	"regexp"
)

type CoordinateHandler struct {
	tracer trace.Tracer
}

type WeatherOutput struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_k"`
}

type Input struct {
	Cep string `json:"cep"`
}

func (h *CoordinateHandler) Get(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.tracer.Start(ctx, "get-location-data")
	defer span.End()
	
	var i Input
	e := json.NewDecoder(r.Body).Decode(&i)
	if e != nil {
		api.WriteJsonResult(w, http.StatusBadRequest, e.Error())
		return
	}

	match, _ := regexp.MatchString("[0-9]{8}", i.Cep)

	if !match {
		api.WriteJsonResult(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	uri := fmt.Sprintf("http://%s/%s", net.JoinHostPort(configs.Cfg.OrchestratorApiHost, configs.Cfg.OrchestratorApiPort), i.Cep)
	client := resty.New()

	req := client.R().
		SetContext(r.Context()).
		SetHeader("Accept", "application/json").
		SetResult(&WeatherOutput{})

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, e := req.Get(uri)

	if e != nil {
		api.WriteJsonResult(w, http.StatusBadRequest, e.Error())
		return
	}

	result := &WeatherOutput{
		City:       resp.Result().(*WeatherOutput).City,
		Celsius:    resp.Result().(*WeatherOutput).Celsius,
		Fahrenheit: resp.Result().(*WeatherOutput).Fahrenheit,
		Kelvin:     resp.Result().(*WeatherOutput).Kelvin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func NewCoordinateHandler(tracer trace.Tracer) *CoordinateHandler {
	return &CoordinateHandler{
		tracer: tracer,
	}
}

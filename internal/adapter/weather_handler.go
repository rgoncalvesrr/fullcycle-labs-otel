package adapter

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/internal/application"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/pkg/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type WeatherHandler struct {
	tracer trace.Tracer
}

func NewWeatherHandler(tracer trace.Tracer) *WeatherHandler {
	return &WeatherHandler{
		tracer: tracer,
	}
}

func (h *WeatherHandler) Get(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := h.tracer.Start(ctx, "get-weather-data")
	defer span.End()

	cep := chi.URLParam(r, "cep")
	coordinateRepository := NewCoordinateRepository(h.tracer)
	weatherRepository := NewWeatherRepository(h.tracer)

	s := application.NewWeatherService(coordinateRepository, weatherRepository)
	output, e := s.GetTemperature(ctx, cep)

	if e != nil {
		switch e {
		case application.ErrCepInvalid, application.ErrCepMalformed:
			api.WriteJsonResult(w, http.StatusUnprocessableEntity, "invalid zipcode")
		case application.ErrCepNotFound:
			api.WriteJsonResult(w, http.StatusNotFound, "can not find zipcode")
		default:
			api.WriteJsonResult(w, http.StatusInternalServerError, e.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(output)
}

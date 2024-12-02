package adapter

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/configs"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/internal/application"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"regexp"
)

type coordinateRepository struct {
	tracer trace.Tracer
}

type Output struct {
	City string `json:"city"`
	Lat  string `json:"lat"`
	Lng  string `json:"lng"`
}

func NewCoordinateRepository(tracer trace.Tracer) application.ICoordinateRepository {
	return &coordinateRepository{
		tracer: tracer,
	}
}

func (c *coordinateRepository) GetByCep(ctx context.Context, cep string) (*application.Coordinate, error) {

	ctx, span := c.tracer.Start(ctx, "get-external-cep")
	defer span.End()

	match, _ := regexp.MatchString("[0-9]{8}", cep)

	if !match {
		return nil, application.ErrCepMalformed
	}

	url := fmt.Sprintf("%s/{cep}", configs.Cfg.CepApiUrl)

	client := resty.New()
	r, e := client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("cep", cep).
		SetResult(&Output{}).
		Get(url)

	if e != nil {
		return nil, e
	}

	switch r.StatusCode() {
	case http.StatusNotFound:
		return nil, application.ErrCepNotFound
	case http.StatusBadRequest:
		return nil, application.ErrCepInvalid
	}

	cord := application.NewCoordinate(r.Result().(*Output).City, r.Result().(*Output).Lat, r.Result().(*Output).Lng)
	return cord, nil
}

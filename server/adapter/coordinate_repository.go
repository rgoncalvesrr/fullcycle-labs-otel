package adapter

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/application"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/configs"
	"net/http"
	"regexp"
)

type coordinateRepository struct {
}

type Output struct {
	City string `json:"city"`
	Lat  string `json:"lat"`
	Lng  string `json:"lng"`
}

func NewCoordinateRepository() application.ICoordinateRepository {
	return &coordinateRepository{}
}

func (c *coordinateRepository) GetByCep(ctx context.Context, cep string) (*application.Coordinate, error) {
	url := fmt.Sprintf("%s/{cep}", configs.Cfg.CepApiUrl)

	match, _ := regexp.MatchString("[0-9]{8}", cep)

	if !match {
		return nil, application.ErrCepMalformed
	}

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

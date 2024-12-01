package application

import "errors"

var (
	ErrCepNotFound        = errors.New("CEP not found")
	ErrCepMalformed       = errors.New("CEP malformed")
	ErrCepInvalid         = errors.New("CEP invalid")
	ErrInvalidTemperature = errors.New("temperature cannot be less than 273.15")
)

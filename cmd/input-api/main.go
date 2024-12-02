package main

import (
	"context"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/configs"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/internal/adapter"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/internal/infra"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/pkg/otel_provider"
	"go.opentelemetry.io/otel"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := otel_provider.InitProvider(configs.Cfg.InputApiOtelServiceName, configs.Cfg.OpenTelemetryCollectorExporterEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shut down Tracer: %s", err)
		}
	}()

	tracer := otel.Tracer("input-api")

	w := infra.NewWebServer(configs.Cfg.InputApiHttpPort)
	w.RegisterHandler("POST /", adapter.NewCoordinateHandler(tracer).Get)

	go func() {
		log.Printf("Iniciando servidor na porta %s", configs.Cfg.InputApiHttpPort)
		w.Start()
	}()

	select {
	case <-signChan:
		log.Println("Shutting down gracefully, Ctrl+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to ohter reason...")
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}

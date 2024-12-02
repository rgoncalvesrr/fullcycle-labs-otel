package otel_provider

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitProvider configura a comunicação com Opentelemetry para ser usado pela aplicação.
func InitProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
	ctx := context.Background()

	// Criação de recurso. Define o nome do serviço que será exibido no tracer (Jaeger, Zipkin)
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Muda o contexto inicial para contexto com tempo limite.
	// Caso não haja comunicação com o coletor em 1 segundo gera erro de comunicação.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Abrindo canal de comunicação com o coletor com contexto de cancelamento por tempo limite.
	conn, err := grpc.DialContext(
		ctx,
		collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Configura o exportar do tracer. Aqui utiliza-se a conexão grpc "otlptracegrpc.WithGRPCConn(conn)".
	// É possível utilizar a conexão do tipo http
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Tipo de exportação que será utilizada. Neste caso será feito com batch (em lote).
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Tipo de amostragem enviada do coletor. Em produção evitar usar AlwaysSample
		sdktrace.WithResource(res),                    // Recurso criado no início da rotina
		sdktrace.WithSpanProcessor(bsp),               // Span processor
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Retornando função de desligamento gracioso
	return tracerProvider.Shutdown, nil
}

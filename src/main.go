package main

import (
	"context"
	"file-storage-server/internal/pkg/config"
	"file-storage-server/internal/pkg/health"
	"file-storage-server/internal/pkg/storage"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		signalTrapped := <-c
		fmt.Println(signalTrapped)
		cancel()
	}()

	cfg := config.NewConfig()

	if err := run(ctx, cfg); err != nil {
		fmt.Println(err)
	}
}

func run(ctx context.Context, cfg *config.Config) (err error) {
	if _, err := os.Stat(cfg.StoragePath); os.IsNotExist(err) {
		os.MkdirAll(cfg.StoragePath, 0644)
	}

	tp := initTracerProvider()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	mp := initMeterProvider()
	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}()

	tracer = tp.Tracer("storageservice")

	r := mux.NewRouter()
	r.HandleFunc("/health", health.HealthGet())
	r.HandleFunc("/files/{file_name}", storage.UploadFileHandler()).Methods(http.MethodPost)
	r.HandleFunc("/files/{file_name}", storage.DeleteFileHandler()).Methods(http.MethodDelete)
	r.HandleFunc("/files", storage.GetAllFilesHandler()).Methods(http.MethodGet)
	http.Handle("/", r)

	s := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: otelhttp.NewHandler(r, "server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)),
	}

	sErrChan := make(chan error, 1)
	go func() {
		err = s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			sErrChan <- err
		}
	}()

	select {
	case <-ctx.Done():
	case sErr := <-sErrChan:
		err = sErr
	}

	if err != nil {
		return
	}

	ctxClose, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = s.Shutdown(ctxClose); err != nil {
		return
	}

	if err == http.ErrServerClosed {
		err = nil
	}
	return
}

func initTracerProvider() *sdktrace.TracerProvider {
	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		log.Fatalf("new otlp trace grpc exporter failed: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func initMeterProvider() *sdkmetric.MeterProvider {
	ctx := context.Background()

	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		log.Fatalf("new otlp metric grpc exporter failed: %v", err)
	}

	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)))
	global.SetMeterProvider(mp)
	return mp
}

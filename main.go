package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()

	shutdown, err := setupOTelSDK(ctx, "test-app", "0.0.1")
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)

	http.HandleFunc("/test", httpHandler)

	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		IdleTimeout:  15 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      nil,
	}
	srv.ListenAndServe()
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.GetTracerProvider().Tracer("tracer")

	_, span := tracer.Start(r.Context(), "hello-span")
	defer span.End()

	meter := otel.GetMeterProvider().Meter("meter")

	counter, _ := meter.Int64Counter("counter")
	counter.Add(context.Background(), 100)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`hello`))
}

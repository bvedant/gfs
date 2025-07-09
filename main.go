package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		srw := &statusResponseWriter{ResponseWriter: w, status: 200}

		next.ServeHTTP(srw, r)

		duration := time.Since(start)

		slog.Info("request",
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", srw.status),
			slog.Duration("duration", duration),
		)
	})
}

func main() {
	dir := flag.String("dir", ".", "directory to serve files from")
	port := flag.Int("port", 8080, "port to listen on")

	flag.Parse()

	fs := http.FileServer(http.Dir(*dir))

	http.Handle("/", loggingMiddleware(fs))

	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	slog.Info("Starting server", slog.String("dir", *dir), slog.String("addr", addr))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ListenAndServe failed", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	<-stop

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}
	slog.Info("Server stopped")
}

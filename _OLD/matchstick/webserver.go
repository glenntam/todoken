package matchstick

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	reqReadTimeout            = 5 * time.Second
	reqWriteTimeout           = 10 * time.Second
	reqIdleTimeout            = 1 * time.Minute
	wsGracefulShutdownTimeout = 30 * time.Second
)

// WebServer is a *http.Server that wraps a safer startup and graceful shutdown.
type WebServer struct {
	*http.Server
}

// NewWebServer ensures sensible timeouts to prevent DoS.
//
// colonPort should include a colon (for example, ":8080").
func NewWebServer(colonPort string, r *Router) *WebServer {
	return &WebServer{
		&http.Server{
			Addr:         colonPort,
			Handler:      r,
			ReadTimeout:  reqReadTimeout,
			WriteTimeout: reqWriteTimeout,
			IdleTimeout:  reqIdleTimeout,
		},
	}
}

// Start runs ListenAndServe in a goroutine and gracefully shuts down on appropriate os signals.
func (ws *WebServer) Start() {
	// Start server in a goroutine
	go func() {
		slog.Info("WebServer starting", "port", ws.Addr)
		if err := ws.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("couldn't start server: %v", err))
		}
	}()

	// Wait for OS interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	slog.Info("shutting down server")

	// Graceful shutdown with 30s timeout
	ctx, cancel := context.WithTimeout(context.Background(), wsGracefulShutdownTimeout)
	defer cancel()

	if err := ws.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}
	slog.Info("server stopped gracefully")
}

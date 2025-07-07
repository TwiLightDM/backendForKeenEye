package main

import (
	"backendForKeenEye/internal/container"
	"backendForKeenEye/internal/router"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Backend for KeenEye
// @version 1.0.0
// @description Backend for KeenEye

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey BasicAuth
// @in header
// @name Authorization
func main() {
	c := container.NewContainer()
	r := router.NewRouter(c)

	srv := &http.Server{
		Addr:    ":" + c.Cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server failed: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Server exited gracefully")
	}

	c.PGClient.Close()
	fmt.Println("Database connection closed.")
}

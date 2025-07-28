package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slilp/go-wallet/internal/config"
	"github.com/slilp/go-wallet/internal/middleware"
	"github.com/slilp/go-wallet/internal/servers"
	"github.com/slilp/go-wallet/internal/servers/api_gen"
)

func main() {
	config.InitConfig()

	db, err := servers.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	if err := servers.InitMigrations(db); err != nil {
		log.Println("Error applying migrations:", err)
	}

	server := servers.NewHttpServer(db)

	if server == nil {
		log.Fatal("server initialization failed")
	}

	r := gin.Default()

	r.Use(middleware.AuthAccessTokenMiddleware)

	api_gen.RegisterHandlers(r, server)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%s", config.Config.AppPort),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
	}()

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

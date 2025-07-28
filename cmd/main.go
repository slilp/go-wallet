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
	"github.com/slilp/go-wallet/internal/restapis"
	"github.com/slilp/go-wallet/internal/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/server"
)

func main() {
	config.InitConfig()

	app := server.NewApiServer()
	httpServer := restapis.NewHttpServer(app)

	r := gin.Default()

	r.Use(middleware.AuthAccessTokenMiddleware)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api_gen.RegisterHandlers(r, &httpServer)

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

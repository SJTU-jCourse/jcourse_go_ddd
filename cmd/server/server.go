package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/app"
	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/interface/web"
)

const (
	SuccessCode         = 0
	SuccessMessage      = "success"
	DefaultPort         = "8080"
	HealthCheckEndpoint = "/health"
	HTTPStatusOK        = 200
)

type Server struct {
	server *http.Server
}

func NewServer(serviceContainer *app.ServiceContainer) *Server {
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(web.CORSMiddleware())

	// Register routes
	web.RegisterRouter(router, serviceContainer)

	// Add health check endpoint
	router.GET(HealthCheckEndpoint, func(c *gin.Context) {
		response := dto.BaseResponse{
			Code: SuccessCode,
			Msg:  SuccessMessage,
		}
		c.JSON(HTTPStatusOK, response)
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting unified server on %s", addr)
	log.Printf("Health check available at http://localhost%s/health", addr)

	// Create HTTP server
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
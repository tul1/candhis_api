package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const ShutdownTimeout = 5 * time.Second

type Server struct {
	port       int
	router     *gin.Engine
	httpServer *http.Server
}

func NewGinServer(publicURL string, port int) (*Server, error) {
	s := &Server{port: port, router: gin.New()}

	gin.SetMode(gin.ReleaseMode)

	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "invalid API path")
	})

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           s.router,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return s, nil
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

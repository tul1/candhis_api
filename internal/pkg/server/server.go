package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const ShutdownTimeout = 5 * time.Second

type Server struct {
	port       int
	router     *gin.Engine
	httpServer *http.Server
}

func NewGinServer(log *logrus.Logger, publicURL string, port int) (*Server, error) {
	s := &Server{port: port, router: gin.New()}

	gin.SetMode(gin.ReleaseMode)

	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "invalid API path")
	})
	s.router.Use(logRequestMiddleware(log))

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

func logRequestMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		var body []byte
		if c.Request.Body != nil {
			var buf bytes.Buffer
			var err error
			tee := io.TeeReader(c.Request.Body, &buf)
			body, err = io.ReadAll(tee)
			if err != nil {
				log.Errorf("fail to read body: %v", err)
			}
			c.Request.Body = io.NopCloser(&buf)
		}

		logWithFields := log.WithFields(
			logrus.Fields{
				"start":     start,
				"path":      c.Request.URL.Path,
				"path_rule": c.FullPath(),
				"query":     c.Request.URL.RawQuery,
				"method":    c.Request.Method,
				"body":      string(body),
			})
		logWithFields.Info("request received")

		// Process request
		c.Next()

		end := time.Now()
		logWithFields.WithFields(logrus.Fields{
			"status_code": c.Writer.Status(),
			"end":         end,
			"latency":     end.Sub(start),
		}).Info("request handled")
	}
}

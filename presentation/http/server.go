package httpserver

import (
	"fmt"
	"net/http"
	"proposal-template/pkg/logger"
	utils "proposal-template/pkg/utils/config"

	"github.com/gin-gonic/gin"
)

type serverConfig struct {
	address string
	port    int
}

var DefaultConfig = serverConfig{
	address: "localhost",
	port:    8080,
}

type HTTPServer struct {
	config utils.Config
	logger logger.ILogger
	router *gin.Engine
}

type Option func(*HTTPServer)
func NewHTTPServer(opts ...Option) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)

	hs := &HTTPServer{
		router: gin.Default(),
	}

	// Apply functional options
	for _, opt := range opts {
		opt(hs)
	}

	// Final setup
	hs.SetupRouter()

	return hs
}

func (s *HTTPServer) Initialize() *HTTPServer {
	s.SetupRouter()
	return s
}

func (s *HTTPServer) SetupRouter() {
	s.logger.Info("Initializing routes...")
	// s.router.Use(middleware.RequestInfoMiddleware(*s.svcCtx))
	s.addRoute(nil, "GET", "/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	v1 := s.router.Group("/api/v1")
	{
		
		s.SetupUserRouter(v1)
	}
}

func (s *HTTPServer) addRoute(group *gin.RouterGroup, method string, path string, handler gin.HandlerFunc) {
	if group == nil {
		s.router.Handle(method, path, handler)
		s.logger.Info(fmt.Sprintf("Route initialized - Method: %s, Path: %s", method, path))
	} else {
		group.Handle(method, path, handler)
		s.logger.Info(fmt.Sprintf("Route initialized - Method: %s, Path: %s", method, group.BasePath()+path))
	}
}
func (s *HTTPServer) Start() error {

	// addr := fmt.Sprintf("%s:%d", DefaultConfig.address, DefaultConfig.port)
	// fmt.Println("logger:", s.logger)
	// s.logger.Info(fmt.Sprintf("Starting HTTP server on address: %v...", addr))
	if err := s.router.Run(); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to start HTTP server: %s", err))
		return err
	}

	return nil
}

// === Optional configuration like logger, system config,.... ===
func WithLogger(logger logger.ILogger) Option {
	return func(s *HTTPServer) {
		s.logger = logger
	}
}

func WithConfig(config utils.Config) Option {
	return func(s *HTTPServer) {
		s.config = config
	}
}

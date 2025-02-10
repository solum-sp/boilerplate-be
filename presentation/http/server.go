package httpserver

import (
	"fmt"
	"net/http"
	"proposal-template/pkg/logger"
	utils "proposal-template/pkg/utils/config"

	"github.com/gin-gonic/gin"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// "github.com/swaggo/files"
)


var DefaultConfig = utils.HttpServerConfig{
	Host: "localhost",
	Port: 8080,
}

type HTTPServer struct {
	config utils.HttpServerConfig
	logger logger.ILogger
	router *gin.Engine
}

type Option func(*HTTPServer)
func NewHTTPServer(opts ...Option) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)

	hs := &HTTPServer{
		config: DefaultConfig,
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
	// s.ServeSwagger()
	s.addRoute(nil, "GET", "/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	v1 := s.router.Group("/api/v1")
	{
		
		s.SetupUserRouter(v1)
	}
}

// func (s *HTTPServer) ServeSwagger() {
// 	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// }

// addRoute adds a route to the HTTP server. If the group parameter is nil, the route is added to the root router.
// Otherwise, the route is added to the given group. The description parameter is optional and is used to provide a description for the route.
func (s *HTTPServer) addRoute(group *gin.RouterGroup, method string, path string, handler gin.HandlerFunc, description ...string) {
	desc := "No description provided" // Default if empty

	if len(description) > 0 {
		desc = description[0] // Use first argument if provided
	}

	if group == nil {
		s.router.Handle(method, path, handler)
		s.logger.Info(fmt.Sprintf("Route initialized - Method: %s, Path: %s, Description: %s", method, path, desc))
	} else {
		group.Handle(method, path, handler)
		s.logger.Info(fmt.Sprintf("Route initialized - Method: %s, Path: %s, Description: %s", method, group.BasePath()+path, desc))
	}
}
func (s *HTTPServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	s.logger.Info(fmt.Sprintf("Starting HTTP server on address: %v...", addr))
	if err := s.router.Run(addr); err != nil {
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

func WithConfig(config utils.HttpServerConfig) Option {
	return func(s *HTTPServer) {
		if config == (utils.HttpServerConfig{}) { // Prevent assigning an empty config
			s.config = DefaultConfig
		} else {
			s.config = config
		}
		fmt.Println("config port:", s.config.Port)
	}
}
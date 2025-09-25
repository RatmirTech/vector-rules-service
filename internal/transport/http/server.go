package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ratmirtech/vector-rules-service/internal/domain"
)

// Server represents the HTTP server
type Server struct {
	echo            *echo.Echo
	ruleHandler     *RuleHandler
	ruleTypeHandler *RuleTypeHandler
}

// NewServer creates a new HTTP server
func NewServer(
	ruleService domain.RuleService,
	ruleTypeService domain.RuleTypeService,
) *Server {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Handlers
	ruleHandler := NewRuleHandler(ruleService)
	ruleTypeHandler := NewRuleTypeHandler(ruleTypeService)

	server := &Server{
		echo:            e,
		ruleHandler:     ruleHandler,
		ruleTypeHandler: ruleTypeHandler,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// Health check
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// API v1 routes
	v1 := s.echo.Group("/api/v1")

	// Rules routes
	v1.POST("/rules", s.ruleHandler.CreateRule)
	v1.GET("/rules/:id", s.ruleHandler.GetRule)
	v1.PUT("/rules/:id", s.ruleHandler.UpdateRule)
	v1.DELETE("/rules/:id", s.ruleHandler.DeleteRule)
	v1.GET("/rules", s.ruleHandler.ListRules)

	// Rule types routes
	v1.POST("/rule-types", s.ruleTypeHandler.CreateRuleType)
	v1.GET("/rule-types/:id", s.ruleTypeHandler.GetRuleType)
	v1.PUT("/rule-types/:id", s.ruleTypeHandler.UpdateRuleType)
	v1.DELETE("/rule-types/:id", s.ruleTypeHandler.DeleteRuleType)
	v1.GET("/rule-types", s.ruleTypeHandler.ListRuleTypes)
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	return s.echo.Start(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	return s.echo.Close()
}
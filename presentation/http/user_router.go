package httpserver

import (
	"proposal-template/presentation/http/handler"
	"github.com/gin-gonic/gin"
)
// SetupUserRouter configures the routes for the User resource.
// It sets up a group (prefix) of routes for the User resource
// and adds a single route, GET /users/:id, which retrieves a
// User by ID.
func (h *HTTPServer) SetupUserRouter(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	userHandler := handler.NewUserHandler()
	h.addRoute(userGroup, "GET", "/:id", userHandler.GetUserById)
}
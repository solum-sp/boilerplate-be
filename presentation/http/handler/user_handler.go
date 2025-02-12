package handler

import (
	"proposal-template/models"
	"proposal-template/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
)

type IUserService interface {
	GetById(id string) (*model.User, error)
}

type UserHandler struct {
	logger logger.ILogger
	UserService IUserService
}

type Option func(*UserHandler)

func NewUserHandler(opts ...Option) *UserHandler {

	var userService IUserService
	container.Resolve(&userService)

	userHandler := &UserHandler{
		UserService: userService,
	}

	for _, opt := range opts {
		opt(userHandler)
	}
	return userHandler
}

func (u *UserHandler) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	data, err := u.UserService.GetById(id)
	if err != nil {
		u.logger.Error("Error getting user by id: " + err.Error())
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(200, gin.H{"data": data})
}

// === optional dependencies ===
func WithLogger(logger logger.ILogger) Option {
	return func(h *UserHandler) {
		h.logger = logger
	}
}

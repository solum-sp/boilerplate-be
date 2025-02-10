package handler

import (

	"proposal-template/models"
	"proposal-template/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
)
type IUserService interface {
	GetUserById(id string) (*model.User, error)
}

type UserHandler struct {
	logger logger.ILogger
	IUserService
}

type Option func(*UserHandler)
func NewUserHandler(opts ...Option) *UserHandler {
	
	var UserService IUserService
	container.Resolve(&UserService)

	userHandler := &UserHandler{
		IUserService: UserService,
	}

	for _, opt := range opts {
		opt(userHandler)
	}
	return userHandler
}


func (u *UserHandler) GetUserById(ctx *gin.Context)  {
	
}

func WithLogger(logger logger.ILogger) Option {
	return func(h *UserHandler) {
		h.logger = logger
	}
}	
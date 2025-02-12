package biz

import (
	"context"
	model "proposal-template/models"
	"proposal-template/pkg/logger"

	// "github.com/golobby/container/v3"
	"gorm.io/gorm"
)

type IUserRepo interface {
	GetByColumn(ctx context.Context, column string, value interface{}) (*model.User, error)
	List(ctx context.Context, paging model.Paging, query *gorm.DB) ([]model.User, error)
	Create(ctx context.Context, user model.User) (uint, error)
}

type UserService struct {
	repo IUserRepo
	logger logger.ILogger
}

type Option func (*UserService)
func NewUserService(repo IUserRepo, opts ...Option) *UserService {

	userService := &UserService{
		repo: repo,
	}

	for _, opt := range opts {
		opt(userService)
	}
	return userService
}

func (s *UserService)  GetById(id string) (*model.User, error) {
	return s.repo.GetByColumn(context.Background(), "id", id)
}

func WithLogger(logger logger.ILogger) Option {
	return func(h *UserService) {
		h.logger = logger
	}
}


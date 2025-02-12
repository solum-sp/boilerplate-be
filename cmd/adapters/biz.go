package adapters

import (
	"fmt"
	"proposal-template/biz"
	"proposal-template/pkg/logger"
	"proposal-template/presentation/http/handler"

	"github.com/golobby/container/v3"
)

func IoCBiz() {
	container.TransientLazy(func() handler.IUserService{
		var (
			logger  logger.ILogger
			userRepo biz.IUserRepo
		)

		err := container.Resolve(&logger)
		if err != nil {
			panic(err)
		}
		err = container.Resolve(&userRepo)
		if err != nil {
			panic(err)
		}
		
		userService := biz.NewUserService(
			userRepo,
			biz.WithLogger(logger),
		)
		fmt.Println("UserService successfully registered in IoC")

		return userService
	})
}
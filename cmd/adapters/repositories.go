package adapters

import (
	"fmt"
	"proposal-template/biz"
	"proposal-template/datalayers/datasources/repositories"

	"github.com/golobby/container/v3"
	"gorm.io/gorm"
)

func IoCRepositories() {
	container.Singleton(func() biz.IUserRepo {
		var (
			db *gorm.DB
		)

		container.Resolve(&db)
		userRepo := repositories.NewUserRepo(db)
		fmt.Println("UserRepo successfully registered in IoC")
		return userRepo
	})
}

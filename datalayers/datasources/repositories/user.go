package repositories

import (
	"database/sql"
	"proposal-template/model"
	"proposal-template/pkg/DAO"
)

type UserRepo struct {
	*dao.GenericDAO[model.User]
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		GenericDAO: dao.NewGenericDAO[model.User](db, "users"),
	}
}
package repositories

import (
	"database/sql"
	"proposal-template/models"
	"proposal-template/pkg/DAO"
)

var usersTableName = "users"
type UserRepo struct {
	*dao.GenericDAO[model.User]
}

// NewUserRepo creates a new UserRepo instance with a GenericDAO for handling
// database operations on the "users" table. It takes a sql.DB connection
// and returns a pointer to the UserRepo.

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		GenericDAO: dao.NewGenericDAO[model.User](db, usersTableName),
	}
}

// === Implement other methods of UserRepo below ==

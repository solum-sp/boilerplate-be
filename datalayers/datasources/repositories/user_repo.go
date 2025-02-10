package repositories

import (
	"database/sql"
	"proposal-template/models"
)

// Define the table name for Users
var usersTableName = "users"

type UserRepo struct {
	*GenericDAO[model.User] 
}

// NewUserRepo creates a new UserRepo instance
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		GenericDAO: NewGenericDAO[model.User](db, usersTableName), 
	}
}
// === Implement other methods of UserRepo below ==

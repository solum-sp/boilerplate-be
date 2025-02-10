package repositories

import (

	"proposal-template/models"

	"gorm.io/gorm"
)

// Define the table name for Users
var usersTableName = "users"

type UserRepo struct {
	*GenericDAO[model.User] 
}

// NewUserRepo creates a new UserRepo instance
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		GenericDAO: NewGenericDAO[model.User](db, usersTableName), 
	}
}
// === Implement other methods of UserRepo below ==

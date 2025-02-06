package model


type User struct {
	BaseModel
	Name          string     `json:"name" db:"name"`
	Email         string     `json:"email" db:"email"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
}



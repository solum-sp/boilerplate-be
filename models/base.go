package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	UUID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
}


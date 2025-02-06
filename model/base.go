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

type Paging struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (p *Paging) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
}

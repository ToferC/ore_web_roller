package models

import (
	"time"

	"github.com/toferc/oneroll"
)

// PowerModel implements a Web model for oneroll.Power
type PowerModel struct {
	ID        int64
	Author    *User
	Power     *oneroll.Power
	Open      bool
	Likes     int
	Slug      string
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}

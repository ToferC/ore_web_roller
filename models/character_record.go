package models

import (
	"time"

	"github.com/toferc/oneroll"
)

type CharacterModel struct {
	ID        int64
	Author    *User
	Character *oneroll.Character
	Open      bool
	Likes     int
	Image     *Image
	Slug      string
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}

type Image struct {
	Id   int
	Path string
}

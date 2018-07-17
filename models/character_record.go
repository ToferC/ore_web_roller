package models

import "github.com/toferc/oneroll"

type CharacterModel struct {
	ID        int64
	Author    User
	Character *oneroll.Character
	Open      bool
	Likes     int
	Slug      string
}

package models

import "github.com/toferc/oneroll"

type PowerModel struct {
	ID     int64
	Author User
	Power  *oneroll.Power
	Open   bool
	Likes  int
	Slug   string
}

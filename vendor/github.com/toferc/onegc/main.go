package main

import (
	"fmt"
	"math/rand"

	"github.com/toferc/onegc/model"
)

func main() {
	fmt.Println("Welcome to OneGC")

	r := rand.New(rand.NewSource(99))

	u := model.User{
		ID:       r.Int63(),
		Email:    "default@gmail.com",
		Password: "12345",
		Profile:  &model.Profile{},
	}

	p := model.Profile{
		Name: "Jane Doe",
		DOB:  "1970-02-23",
	}

	u.Profile = &p

	o := model.Organization{
		ID:   0,
		Name: "Canada Revenue Agency",
	}

	s := model.Service{
		Name:    "Taxes",
		Actions: map[string]*model.Action{},
	}

	a := model.Action{
		Name: "Collect",
	}
}

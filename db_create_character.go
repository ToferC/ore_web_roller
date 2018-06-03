package main

import (
	"fmt"
	"strconv"

	"github.com/fatih/structs"
	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
)

// CreateCharacter takes terminal user input and saves to DB
func CreateCharacter(db *pg.DB) *oneroll.Character {

	name := UserQuery("What is the character's name? ")

	c := oneroll.NewWTCharacter(name)

	m := structs.Map(c)
	m["Name"] = name

	// Add statistics

	fmt.Println("\nAdding stats and skills. We'll start with normal dice and you can update to add special abilities later.")

	statistics := []*oneroll.Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

	for _, s := range statistics {
		answer := UserQuery("Input normal die value for " + s.Name + ": ")
		num, _ := strconv.Atoi(answer)

		s.Dice.Normal = num
	}

	c.BaseWill = c.Command.Dice.Normal + c.Charm.Dice.Normal
	c.Willpower = c.BaseWill

	for k, v := range c.Skills {
		str := fmt.Sprintf("Input value for %s (%s): ", k, v.LinkStat.Name)
		answer := UserQuery(str)
		num, _ := strconv.Atoi(answer)

		c.Skills[k].Dice.Normal = num
	}

	fmt.Println(c)

	// Save character
	err := SaveCharacter(db, c)
	if err != nil {
		panic(err)
	}

	return c
}

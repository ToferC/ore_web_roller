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

	fmt.Println("\nAdding stats and skills.")

	statistics := []*oneroll.Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

	fmt.Println("Enter normal die values for:")

	for _, s := range statistics {
		answer := UserQuery("\n" + s.Name + ": ")
		num, _ := strconv.Atoi(answer)

		s.Dice.Normal = num
		for k, v := range c.Skills {
			if v.LinkStat.Name == s.Name {
				str := fmt.Sprintf("-- %s: ", k)
				answer := UserQuery(str)
				num, _ := strconv.Atoi(answer)

				c.Skills[k].Dice.Normal = num
			}
		}
	}

	c.BaseWill = c.Command.Dice.Normal + c.Charm.Dice.Normal
	c.Willpower = c.BaseWill

	fmt.Println(c)

	// Save character
	err := SaveCharacter(db, c)
	if err != nil {
		panic(err)
	}

	return c
}

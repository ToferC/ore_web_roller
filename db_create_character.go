package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
)

// CreateCharacter takes terminal user input and saves to DB
func CreateCharacter(db *pg.DB) *oneroll.Character {

	question := bufio.NewReader(os.Stdin)
	fmt.Print("What is your name? ")
	r, _ := question.ReadString('\n')

	name := strings.Trim(r, " \n")

	c := oneroll.NewCharacter(name)

	m := structs.Map(c)
	m["Name"] = name

	// Add statistics

	statistics := []*oneroll.Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

	for _, s := range statistics {
		q := bufio.NewReader(os.Stdin)
		fmt.Print("Input value for ", s.Name, ": ")
		r, _ = q.ReadString('\n')

		answer := strings.Trim(r, " \n")
		num, _ := strconv.Atoi(answer)

		s.Dice.Normal = num
	}

	c.BaseWill = c.Command.Dice.Normal + c.Charm.Dice.Normal
	c.Willpower = c.BaseWill

	for k := range c.Skills {
		question = bufio.NewReader(os.Stdin)
		fmt.Print("Input value for ", k, ": ")
		r, _ = question.ReadString('\n')

		answer := strings.Trim(r, " \n")
		num, _ := strconv.Atoi(answer)

		c.Skills[k].Dice.Normal = num
	}

	c.Display()

	// Save character
	err := SaveCharacter(db, c)
	if err != nil {
		panic(err)
	}

	return c
}

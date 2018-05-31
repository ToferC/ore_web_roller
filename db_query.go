package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
)

// Query and return a Character from DB
func Query(db *pg.DB) {

	// List all charcters in DB
	err := ListCharacters(db)
	if err != nil {
		panic(err)
	}

	// Get user input on which character to load
	question := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your character's name to load: ")
	q, _ := question.ReadString('\n')

	name := strings.Trim(q, " \n")

	c := LoadCharacter(db, name)
	c.Display()

	s := oneroll.FormSkillDieString(c.Skills["Athletics"], 1)

	fmt.Printf("Rolling Athletics (Body+Athletics) for %s\n", c.Name)
	r := oneroll.Roll{
		Actor:  c,
		Action: "act",
	}

	r.Resolve(s)

	fmt.Println("Rolling!")
	fmt.Println(r)

	fmt.Println(c)

}

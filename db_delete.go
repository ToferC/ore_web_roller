package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-pg/pg"
)

// Delete removes a Character from the DB
func Delete(db *pg.DB) {

	// List all characters in DB
	err := ListCharacters(db)
	if err != nil {
		panic(err)
	}

	// Get user input on which character to load
	question := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your character's name to delete: ")
	q, _ := question.ReadString('\n')

	name := strings.Trim(q, " \n")

	c := LoadCharacter(db, name)
	c.Display()

	question = bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure you want to delete ", c.Name, " ? (Y/N)")
	q, _ = question.ReadString('\n')

	response := strings.Trim(q, " \n")

	if response == "Y" || response == "y" {
		err = DeleteCharacter(db, c.ID)
		if err != nil {
			panic(err)
		}
		fmt.Println("Deleted.")
	} else {
		fmt.Println("Delete aborted.")
	}
}

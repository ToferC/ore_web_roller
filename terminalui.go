package main

import (
	"fmt"

	"github.com/go-pg/pg"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/toferc/ore_web/terminal"
)

// Terminal runs the basic functions from a terminal input
func Terminal(db *pg.DB) {

MainLoop:
	for true {
		fmt.Println("\nWelcome to the ORE Terminal")
		fmt.Println("\nFrom here, you can query, create, update or delete ORE characters.")
		fmt.Println("\nPlease select your action.")

		input := terminal.UserQuery(`
  1: Query the Database
  2: Create a Character
  3: Update a Character
  4: Delete a Character

  5: Start a Conflict (Coming Soon)

  Selection: `)

		switch input {
		case "1":
			terminal.Query(db)
		case "2":
			terminal.CreateCharacter(db)
		case "3":
			terminal.Update(db)
		case "4":
			terminal.Delete(db)
		default:
			fmt.Println("Invalid input. Exiting.")
			break MainLoop
		}
	}
}

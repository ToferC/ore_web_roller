package main

import (
	"fmt"

	"github.com/go-pg/pg"
)

// Delete removes a Character from the DB
func Delete(db *pg.DB) {

	c, err := GetCharacter(db)
	if err != nil {
		panic(err)
	}

	response := UserQuery("Are you sure you want to delete " + c.Name + " ? (Y/N)")

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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-pg/pg"
)

func init() {
	os.Setenv("DBUser", "chris")
	os.Setenv("DBPass", "12345")
	os.Setenv("DBName", "ore_engine")
}

// Terminal runs the basic functions from a terminal input
func Terminal() {

	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DBUser"),
		Password: os.Getenv("DBPass"),
		Database: os.Getenv("DBName"),
	})

	defer db.Close()

	fmt.Println("Welcome to the ORE Terminal")
	fmt.Println("From here, you can query, create or delete ORE characters.")
	fmt.Println("Please select your action.")

	input := UserQuery(
		`
  1: Query the Database
  2: Create a Character
  3: Delete a Character

  Selection: `)

	switch input {
	case "1":
		Query(db)
	case "2":
		CreateCharacter(db)
	case "3":
		Delete(db)
	default:
		fmt.Println("Invalid input. Exiting.")
	}
}

// UserQuery creates and question and returns the User's input as a string
func UserQuery(q string) string {
	question := bufio.NewReader(os.Stdin)
	fmt.Print(q)
	r, _ := question.ReadString('\n')

	input := strings.Trim(r, " \n")

	return input
}
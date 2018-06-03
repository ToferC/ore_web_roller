package main

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/toferc/oneroll"
)

// SaveCharacter saves a Character to the DB
func SaveCharacter(db *pg.DB, c *oneroll.Character) error {
	// Save character in Database
	_, err := db.Model(c).
		OnConflict("(id) DO UPDATE").
		Set("name = ?name").
		Insert(c)
	if err != nil {
		panic(err)
	}
	return err
}

func UpdateCharacter(db *pg.DB, c *oneroll.Character) error {
	err := db.Update(c)
	if err != nil {
		panic(err)
	}
	return err
}

// InitDB initializes the DB Schema
func InitDB(db *pg.DB) error {
	err := createSchema(db)
	if err != nil {
		panic(err)
	}
	return err
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*oneroll.Character)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ListCharacters queries Character names and add to slice
func ListCharacters(db *pg.DB) ([]string, error) {
	var chars []string

	err := db.Model((*oneroll.Character)(nil)).
		Column("name").
		Order("id ASC").
		Select(&chars)

	if err != nil {
		return []string{}, err
	}

	// Print names and PK
	for i, n := range chars {
		fmt.Println(i, n)
	}
	return chars, nil
}

// GetCharacter lists all characters in DB and asks the user to select one
func GetCharacter(db *pg.DB) (*oneroll.Character, error) {

	var name string
	var err error

	// Select character loop
SelectCharacterLoop:
	for true {
		// List all charcters in DB
		list, err := ListCharacters(db)
		if err != nil {
			panic(err)
		}

		// Get user input on which character to load
		name = UserQuery("Enter your character's name to load or hit Enter to quit: ")

		if name == "" {
			fmt.Println("Exiting.")
			break SelectCharacterLoop
		}

		validCharacter := true

		for _, n := range list {
			if name == n {
				validCharacter = true
				break
			}
			validCharacter = false
		}

		if validCharacter == false {
			fmt.Println("Not a valid character. Try again.")
		} else {
			break
		}
	}
	c := LoadCharacter(db, name)
	fmt.Println(c)
	return c, err
}

// LoadCharacter loads a single character from the DB by name
func LoadCharacter(db *pg.DB, name string) *oneroll.Character {
	// Select user by Primary Key
	char := new(oneroll.Character)
	err := db.Model(char).
		Where("Name = ?", name).
		Limit(1).
		Select()

	if err != nil {
		panic(err)
	}

	fmt.Println("Character loaded From DB")
	return char
}

// DeleteCharacter deletes a single character from DB by ID
func DeleteCharacter(db *pg.DB, pk int64) error {

	char := oneroll.Character{ID: pk}

	fmt.Println("Deleting character...")

	err := db.Delete(&char)

	return err
}

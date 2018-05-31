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
func ListCharacters(db *pg.DB) error {
	var chars []string

	err := db.Model((*oneroll.Character)(nil)).
		Column("name").
		Order("id ASC").
		Select(&chars)

	if err != nil {
		return err
	}

	// Print names and PK
	for i, n := range chars {
		fmt.Println(i, n)
	}
	return nil
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

package main

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/toferc/oneroll"
)

func main() {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "devpass",
	})

	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	c := oneroll.NewCharacter("Baron")
	c.Display()

	actingSkill := c.Skills["Athletics"]

	actingSkill.Dice.Spray = 1
	actingSkill.Dice.Normal = 3

	err = db.Insert(c)
	if err != nil {
		panic(err)
	}

	// Select user by Primary Key
	char := &oneroll.Character{Id: c.Id}
	err = db.Select(char)
	if err != nil {
		panic(err)
	}

	fmt.Println("From DB")
	char.Display()

}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*oneroll.Character)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

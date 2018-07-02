package database

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/toferc/oneroll"
)

// InitDB initializes the DB Schema
func InitDB(db *pg.DB) error {
	err := createSchema(db)
	if err != nil {
		panic(err)
	}
	return err
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		(*oneroll.Character)(nil),
		(*oneroll.Power)(nil)} {
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

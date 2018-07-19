package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/toferc/ore_web_roller/models"
)

// SaveUser saves a User to the DB
func SaveUser(db *pg.DB, u *models.User) error {

	// Save User in Database
	_, err := db.Model(u).
		OnConflict("(id) DO UPDATE").
		Set("name = ?name").
		Insert(u)
	if err != nil {
		panic(err)
	}
	return err
}

//UpdateUser updates user info
func UpdateUser(db *pg.DB, u *models.User) error {

	err := db.Update(u)
	if err != nil {
		panic(err)
	}
	return err
}

// ListUsers queries User names and add to slice
func ListUsers(db *pg.DB) ([]*models.User, error) {
	var users []*models.User

	_, err := db.Query(&users, `SELECT * FROM Users`)

	if err != nil {
		panic(err)
	}

	// Print names and PK
	for i, n := range users {
		fmt.Println(i, n.Name)
	}
	return users, nil
}

// PKLoadUser loads a single User from the DB by pk
func PKLoadUser(db *pg.DB, pk int64) (*models.User, error) {
	// Select user by Primary Key
	user := &models.User{ID: pk}
	err := db.Select(user)

	if err != nil {
		return &models.User{Name: "New"}, err
	}

	fmt.Println("User loaded From DB")
	return user, nil
}

// DeleteUser deletes a single User from DB by ID
func DeleteUser(db *pg.DB, pk int64) error {

	user := models.User{ID: pk}

	fmt.Println("Deleting User...")

	err := db.Delete(&user)

	return err
}

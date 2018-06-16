package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web/database"
)

func IndexHandler(w http.ResponseWriter, req *http.Request) {

	characters, err := database.ListCharacters(db)
	if err != nil {
		panic(err)
	}

	Render(w, "templates/index.html", characters)
}

// CharacterHandler renders a character in a Web page
func CharacterHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {

		// Render page

		name := req.URL.Path[len("/view/"):]

		if len(name) == 0 {
			name = "Player"
		}

		c, err := database.LoadCharacter(db, name)
		if err != nil {
			fmt.Println(err)
		}

		Render(w, "templates/character.html", c)

	} else {

		// Parse Form and redirect

	}

}

// NewCharacterHandler renders a character in a Web page
func NewCharacterHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {

		name := req.URL.Path[len("/new/"):]

		if len(name) == 0 {
			name = "Player"
		}

		c := oneroll.NewWTCharacter(name)

		// Render page
		Render(w, "templates/characterform.html", c)

	} else {

		c := oneroll.NewWTCharacter("Default")

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		c.Name = req.FormValue("Name")

		for _, st := range c.StatMap {
			c.Statistics[st].Dice.Normal, _ = strconv.Atoi(req.FormValue(st))
		}

		for _, sk := range c.Skills {
			sk.Dice.Normal, _ = strconv.Atoi(req.FormValue(sk.Name))
		}

		err = database.SaveCharacter(db, c)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)
		http.Redirect(w, req, "/character/"+c.Name, http.StatusSeeOther)
	}
}

// ModifyCharacterHandler renders a character in a Web page
func ModifyCharacterHandler(w http.ResponseWriter, req *http.Request) {

	name := req.URL.Path[len("/modify/"):]

	if len(name) == 0 {
		http.Redirect(w, req, "/new/NewCharacter", http.StatusSeeOther)
	}

	c, err := database.LoadCharacter(db, name)
	if err != nil {
		fmt.Println(err)
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/characterform.html", c)

	} else {

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		c.Name = req.FormValue("Name")

		for _, st := range c.StatMap {
			c.Statistics[st].Dice.Normal, _ = strconv.Atoi(req.FormValue(st))
		}

		for _, sk := range c.Skills {
			sk.Dice.Normal, _ = strconv.Atoi(req.FormValue(sk.Name))
		}

		err = database.UpdateCharacter(db, c)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)
		http.Redirect(w, req, "/view/"+c.Name, http.StatusSeeOther)
	}
}

// DeleteCharacterHandler renders a character in a Web page
func DeleteCharacterHandler(w http.ResponseWriter, req *http.Request) {

	pk := req.URL.Path[len("/delete/"):]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_character.html", c)

	} else {

		database.DeleteCharacter(db, c.ID)

		fmt.Println("Deleted ", c.Name)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

package main

import (
	"fmt"
	"net/http"

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

		name := req.URL.Path[len("/character/"):]

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

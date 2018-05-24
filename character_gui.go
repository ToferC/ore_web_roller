package main

import (
	"net/http"

	"github.com/toferc/oneroll"
)

// CharacterHandler renders a character in a Web page
func CharacterHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {

		// Render page

		name := req.URL.Path[len("/character/"):]

		if len(name) == 0 {
			name = "Player"
		}

		c := oneroll.NewCharacter(name)

		c.BaseWill = c.Command.Dice.Normal + c.Charm.Dice.Normal
		c.Willpower = c.BaseWill

		render(w, "templates/character.html", c)

	} else {

		// Parse Form and redirect

	}

}

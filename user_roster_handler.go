package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/ore_web_roller/database"
	"github.com/toferc/ore_web_roller/models"
)

func UserCharacterRosterHandler(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		Render(w, "templates/login.html", nil)
		return
		// in case of error
	}

	// Prep for user authentication
	username := ""

	u := session.Values["username"]

	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		username = user
	}

	characters, err := database.ListUserCharacterModels(db, username)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser:     username,
		CharacterModels: characters,
	}

	Render(w, "templates/user_roster.html", wc)
}

func AddToUserRosterHandler(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		Render(w, "templates/login.html", nil)
		return
		// in case of error
	}

	// Prep for user authentication
	username := ""

	u := session.Values["username"]

	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		username = user
	}

	author := database.LoadUser(db, username)

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	// Load CharacterModel
	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	newCharacterModel := models.CharacterModel{
		Author:    author,
		Character: cm.Character,
		Open:      false,
	}

	err = database.SaveCharacterModel(db, &newCharacterModel)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Saved")
	}

	url := fmt.Sprintf("/view_character/%d", newCharacterModel.ID)

	http.Redirect(w, req, url, 302)
}

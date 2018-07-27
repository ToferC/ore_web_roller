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

//BetaPowerIndexHandler displays a list of PowerModels in a WebCharacter View
func BetaPowerIndexHandler(w http.ResponseWriter, req *http.Request) {

	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		// in case of error
	}

	// Prep for user authentication
	username := ""

	// Get session User
	u := session.Values["username"]

	// Type assertation
	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		username = user
	}

	pows, err := database.ListPowers(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser: username,
		Powers:      pows,
	}

	Render(w, "templates/beta_index_powers.html", wc)
}

func ConvertToOpenPowerArchive(w http.ResponseWriter, req *http.Request) {

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
	p, err := database.PKLoadPower(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	newPowerModel := models.PowerModel{
		Author: author,
		Power:  p,
		Open:   true,
	}

	err = database.SavePowerModel(db, &newPowerModel)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Saved")
		_ = database.DeletePower(db, int64(id))
	}

	url := fmt.Sprintf("/view_power/%d", newPowerModel.ID)

	http.Redirect(w, req, url, 302)
}

func BetaCharacterIndexHandler(w http.ResponseWriter, req *http.Request) {

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

	characters, err := database.ListCharacters(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		SessionUser: username,
		Characters:  characters,
	}

	Render(w, "templates/beta_roster.html", wc)
}

func ConvertToCharacterModel(w http.ResponseWriter, req *http.Request) {

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
	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	newCharacterModel := models.CharacterModel{
		Author:    author,
		Character: c,
		Open:      true,
	}

	err = database.SaveCharacterModel(db, &newCharacterModel)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Saved")
		_ = database.DeleteCharacter(db, int64(c.ID))
	}

	url := fmt.Sprintf("/view_character/%d", newCharacterModel.ID)

	http.Redirect(w, req, url, 302)
}

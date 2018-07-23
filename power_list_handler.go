package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// PowerListHandler renders a character in a Web page
func PowerListHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
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

	// Validate that User == Author
	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	c := cm.Character

	powers, err := database.ListPowers(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		Powers:         powers,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/add_power_from_list.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		if c.Powers == nil {
			c.Powers = map[string]*oneroll.Power{}
		}

		pName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		if pName != "" {
			p := powers[oneroll.ToSnakeCase(pName)]

			fmt.Println(p)

			p.Dice.Normal = nd
			p.Dice.Hard = hd
			p.Dice.Wiggle = wd

			p.Slug = oneroll.ToSnakeCase(pName)

			oneroll.UpdateCost(&p)

			c.Powers[p.Slug] = &p
		}
	}

	err = database.UpdateCharacterModel(db, cm)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Saved")
	}

	fmt.Println(c)

	url := fmt.Sprintf("/view_character/%d", cm.ID)

	http.Redirect(w, req, url, http.StatusSeeOther)
}

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
	"github.com/toferc/ore_web_roller/models"
)

// ModifyAdvantageHandler renders a character in a Web page
func ModifyAdvantageHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	um := &models.User{}
	username := ""

	// Get session User
	u := session.Values["username"]

	// Type assertation
	if user, ok := u.(string); !ok {
		um.UserName = ""
	} else {
		fmt.Println(user)
		username = user
	}

	fmt.Println(um)

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

	// Assign additional empty Sources to populate form
	if len(c.Advantages) < 10 {
		for i := len(c.Advantages); i < 10; i++ {
			tempA := oneroll.Advantage{
				Name:  "",
				Level: 1,
			}
			c.Advantages = append(c.Advantages, &tempA)
		}
	}

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		SessionUser:    username,
		Advantages:     oneroll.Advantages,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/add_advantages.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		c.Advantages = []*oneroll.Advantage{}

		for s := 1; s < 11; s++ {
			aName := req.FormValue(fmt.Sprintf("Advantage-%d-Name", s))

			aInfo := req.FormValue(fmt.Sprintf("Advantage-%d-Info", s))

			if aName != "" {
				a := oneroll.Advantages[aName]
				l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Advantage-%d-Level", s)))
				if err != nil {
					l = 1
				}
				a.Level = l
				a.Info = aInfo
				c.Advantages = append(c.Advantages, &a)
			}
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

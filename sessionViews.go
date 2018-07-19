package main

import (
	"log"
	"net/http"

	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/ore_web_roller/database"
	"github.com/toferc/ore_web_roller/models"
)

//LogoutFunc Implements the logout functionality
//Will delete the session information from the cookie Store
func LogoutFunc(w http.ResponseWriter, req *http.Request) {
	session, err := sessions.Store.Get(req, "session")
	if err == nil {
		if session.Values["loggedin"] != false {
			session.Values["loggedin"] = "false"
			session.Save(req, w)
		}
	}
	http.Redirect(w, req, "/", 302)
	// Redirect to main page
}

func LoginFunc(w http.ResponseWriter, req *http.Request) {
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		// in case of error
	} else {
		isLoggedIn := session.Values["loggedin"]
		if isLoggedIn != true {
			if req.Method == "POST" {
				// first get user object

				if req.FormValue("password") == "secret" && req.FormValue("username") == "user" {
					session.Values["loggedin"] = "true"
					session.Save(req, w)
					http.Redirect(w, req, "/", 302)
					return
				} else if req.Method == "GET" {
					Render(w, "templates/index_characters.html", nil)
				}
			}
		}
	}
}

func SignUpFunc(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()

		username := req.Form.Get("username")
		password := req.Form.Get("password")
		email := req.Form.Get("email")

		log.Println(username, password, email)

		u := models.User{
			Name:     username,
			Password: password,
			Email:    email,
		}

		err := database.SaveUser(db, &u)
		if err != nil {
			http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
		} else {
			Render(w, "/templates/index_characters.html", 302)
		}
	}
}

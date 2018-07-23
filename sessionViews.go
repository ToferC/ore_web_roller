package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/ore_web_roller/database"
	"github.com/toferc/ore_web_roller/models"
)

func UserIndexHandler(w http.ResponseWriter, req *http.Request) {

	users, err := database.ListUsers(db)
	if err != nil {
		panic(err)
	}

	Render(w, "templates/index_users.html", users)
}

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

//LoginFunc implements the login functionality, will add a cookie to cookie Store
//to manage authentication
func LoginFunc(w http.ResponseWriter, req *http.Request) {
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		Render(w, "templates/login.html", nil)
		return
		// in case of error
	}

	switch req.Method {
	case "GET":
		Render(w, "templates/login.html", nil)
	case "POST":
		log.Print("Inside POST")
		req.ParseForm()
		username := req.Form.Get("username")
		password := req.Form.Get("password")

		if (username != "" && password != "") && database.ValidUser(db, username, password) {
			session.Values["loggedin"] = "true"
			session.Values["username"] = username
			session.Save(req, w)
			log.Print("user ", username, " is authenticated")
			fmt.Println(session.Values)
			http.Redirect(w, req, "/", 302)
		} else {
			log.Print("Invalid user " + username)
			Render(w, "templates/login.html", nil)
		}
	}
}

//SignUpFunc implements sign-up functionality
func SignUpFunc(w http.ResponseWriter, req *http.Request) {
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		Render(w, "templates/login.html", nil)
		return
		// in case of error
	}

	if req.Method == "POST" {
		req.ParseForm()

		username := req.Form.Get("username")
		password := req.Form.Get("password")
		email := req.Form.Get("email")

		log.Println(username, password, email)

		u := models.User{
			UserName: username,
			Password: password,
			Email:    email,
		}

		err := database.SaveUser(db, &u)
		if err != nil {
			http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
		} else {
			//Log in user and continue
			session.Values["loggedin"] = "true"
			session.Values["username"] = username
			session.Save(req, w)
			log.Print("user ", username, " is authenticated")
			fmt.Println(session.Values)
			http.Redirect(w, req, "/", 302)
		}
	} else if req.Method == "GET" {
		Render(w, "templates/signup.html", nil)
	}
}

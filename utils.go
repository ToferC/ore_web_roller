package main

import (
	"fmt"

	"github.com/gorilla/sessions"
)

type Element interface {
	ChooseDiePool() error
}

func ChooseDice(e Element) error {
	err := e.ChooseDiePool()
	if err != nil {
		return err
	}
	return nil
}

func getUserSessionValues(s *sessions.Session) map[string]string {

	sessionMap := map[string]string{
		"username": "",
		"loggedin": "false",
		"isAdmin":  "false",
	}

	// Prep for user authentication

	u := s.Values["username"]
	l := s.Values["loggedin"]
	a := s.Values["admin"]

	// Type assertation
	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		sessionMap["username"] = user
	}

	// Type assertation
	if loggin, ok := l.(string); !ok {
	} else {
		fmt.Println(loggin)
		sessionMap["loggedin"] = loggin
	}

	// Type assertation
	if admin, ok := a.(string); !ok {
	} else {
		fmt.Println(admin)
		sessionMap["isAdmin"] = admin
	}
	return sessionMap
}

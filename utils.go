package main

import "net/http"

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

func wrapHandler(
	handler func(w http.ResponseWriter, req *http.Request),
) func(w http.ResponseWriter, req *http.Request) {

	h := func(w http.ResponseWriter, req *http.Request) {
		if !userIsAuthorized(req) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, req)
	}
	return h
}

// https://www.activestate.com/blog/2017/04/creating-web-app-using-golang-gorilla-web-toolkit
func userIsAuthorized(req *http.Request) bool {
	userID := r.Header.Get("X-HashText-User-ID")
	if userID == "" {
		return false
	}

	var found bool
	err := db.QueryRow('SELECT 1 FROM "user" WHERE user_id = $1', userID).Scan(&found)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Printf("Query to look up user failed: %v", err)
		return false
	}
	return found
}

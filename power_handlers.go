package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// AddPowerHandler renders a character in a Web page
func AddPowerHandler(w http.ResponseWriter, req *http.Request) {

	ch := req.URL.Path[len("/add_power/"):]

	id, err := strconv.Atoi(ch)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	if len(c.Powers) == 0 {
		c.Powers = map[string]*oneroll.Power{}
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/form_power.html", c)

	} else { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		pName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		q1name := req.FormValue("Q1-Name")

		q1, q2, q3 := new(*oneroll.Quality)

		if q1name != "" {
			l1, err := strconv.Atoi(req.FormValue("Q1-Level"))
			q1 := &oneroll.Quality{
				Type:  req.FormValue("Q1-Type"),
				Level: l1,
				Name:  req.FormValue("Q1-Name"),
			}
		}

		q2name := req.FormValue("Q2-Name")

		if q2name != "" {
			l2, err := strconv.Atoi(req.FormValue("Q2-Level"))
			q2 := &oneroll.Quality{
				Type:  req.FormValue("Q2-Type"),
				Level: l2,
				Name:  req.FormValue("Q2-Name"),
			}
		}

		q3name := req.FormValue("Q3-Name")

		if q3name != "" {
			l3, err := strconv.Atoi(req.FormValue("Q3-Level"))
			q3 := &oneroll.Quality{
				Type:  req.FormValue("Q3-Type"),
				Level: l3,
				Name:  req.FormValue("Q3-Name"),
			}
		}
		p := oneroll.Power{
			Name: pName,
			Dice: &oneroll.DiePool{
				Normal: nd,
				Hard:   hd,
				Wiggle: wd,
			},
			Effect:    req.FormValue("Effect"),
			Qualities: []*oneroll.Quality{},
		}

		c.Powers[pName] = &p

		err = database.UpdateCharacter(db, c)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)
		http.Redirect(w, req, "/new_power/"+string(c.ID), http.StatusSeeOther)
	}
}

// ModifyPowerHandler renders a character in a Web page
func ModifyPowerHandler(w http.ResponseWriter, req *http.Request) {

	ch := req.URL.Path[len("/modify_power/"):]

	id, err := strconv.Atoi(ch)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	if len(c.Powers) == 0 {
		c.Powers = map[string]*oneroll.Power{}
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/form_power.html", c)

	} else { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		pName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		p := oneroll.Power{
			Name: pName,
			Dice: &oneroll.DiePool{
				Normal: nd,
				Hard:   hd,
				Wiggle: wd,
			},
			Effect: req.FormValue("Effect"),
		}

		c.Powers[pName] = &p

		err = database.UpdateCharacter(db, c)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)
		http.Redirect(w, req, "/new_power/"+string(c.ID), http.StatusSeeOther)
	}
}

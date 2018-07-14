package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// PowerListHandler renders a character in a Web page
func PowerListHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	powers, err := database.ListPowers(db)
	if err != nil {
		panic(err)
	}

	wc := WebChar{
		Character: c,
		Powers:    powers,
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

		if c.Powers != nil {
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

	err = database.UpdateCharacter(db, c)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Saved")
	}

	fmt.Println(c)

	url := fmt.Sprintf("/view_character/%d", c.ID)

	http.Redirect(w, req, url, http.StatusSeeOther)
}

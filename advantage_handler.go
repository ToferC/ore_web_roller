package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// ModifyAdvantageHandler renders a character in a Web page
func ModifyAdvantageHandler(w http.ResponseWriter, req *http.Request) {

	pk := req.URL.Path[len("/add_advantages/"):]

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
		Character:  c,
		Advantages: oneroll.Advantages,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/add_advantages.html", wc)

	} else {

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
				c.Advantages = append(c.Advantages, a)
			}
		}
	}

	err = database.UpdateCharacter(db, c)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Saved")
	}

	fmt.Println(c)
	http.Redirect(w, req, "/view/"+string(c.ID), http.StatusSeeOther)
}

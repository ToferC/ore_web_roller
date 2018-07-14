package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// AddSkillHandler renders a character in a Web page
func AddSkillHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pk := vars["id"]
	s := vars["skill"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Assign basic HyperSkill
	stat := c.Statistics[s]

	skill := &oneroll.Skill{
		Quality: &oneroll.Quality{
			Type: "",
		},
		LinkStat: stat,
		Dice: &oneroll.DiePool{
			Normal: 0,
			Hard:   0,
			Wiggle: 0,
		},
		ReqSpec:        false,
		Specialization: "",
	}

	wc := WebChar{
		Character: c,
		Skill:     skill,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_skill.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		sName := req.FormValue("Name")

		sQuality := req.FormValue("Quality")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))

		skill = new(oneroll.Skill)

		skill.Quality = &oneroll.Quality{
			Type: sQuality,
		}

		skill.Name = sName

		skill.LinkStat = stat

		skill.Dice = &oneroll.DiePool{
			Normal: nd,
		}

		if req.FormValue("Free") != "" {
			skill.Free = true
		}

		if req.FormValue("ReqSpec") == "Yes" {
			skill.ReqSpec = true
			skill.Specialization = req.FormValue("Specialization")
		}

		c.Skills[sName] = skill

		fmt.Println(c)

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
}

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// AddSkillHandler renders a character in a Web page
func AddSkillHandler(w http.ResponseWriter, req *http.Request) {

	s := req.URL.Path[len("/add_skill/"):]

	sSlice := strings.Split(s, "/")

	ch, s := sSlice[0], sSlice[1]

	id, err := strconv.Atoi(ch)
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

	} else { // POST

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
		http.Redirect(w, req, "/view/"+string(c.ID), http.StatusSeeOther)
	}
}

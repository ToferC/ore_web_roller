package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// AddHyperSkillHandler renders a character in a Web page
func AddHyperSkillHandler(w http.ResponseWriter, req *http.Request) {

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
	skill := c.Skills[s]

	skill.HyperSkill = &oneroll.HyperSkill{
		Dice: &oneroll.DiePool{
			Normal: 0,
			Hard:   0,
			Wiggle: 0,
		},
		Effect: "",
	}

	hs := skill.HyperSkill

	hs.Qualities = []*oneroll.Quality{
		&oneroll.Quality{
			Type:  skill.Quality.Type,
			Level: 0,
			Capacities: []*oneroll.Capacity{
				&oneroll.Capacity{
					Type: "Self",
				},
			},
		},
	}

	// Assign additional empty Qualities to populate form
	if len(hs.Qualities) < 4 {
		for i := len(hs.Qualities); i < 4; i++ {
			tempQ := oneroll.NewQuality("")
			hs.Qualities = append(hs.Qualities, tempQ)
		}
	} else {
		// Always create at least 2 Qualities
		for i := 0; i < 2; i++ {
			tempQ := oneroll.NewQuality("")
			hs.Qualities = append(hs.Qualities, tempQ)
		}
	}

	// Assign additional empty Capacities to populate form
	for _, q := range hs.Qualities {
		if len(q.Capacities) < 4 {
			for i := len(q.Capacities); i < 4; i++ {
				tempC := oneroll.Capacity{
					Type: "",
				}
				q.Capacities = append(q.Capacities, &tempC)
			}
		}
		if len(q.Modifiers) < 8 {
			for i := len(q.Modifiers); i < 8; i++ {
				tempM := oneroll.NewModifier("")
				q.Modifiers = append(q.Modifiers, tempM)
			}
		}
	}

	wc := WebChar{
		Character: c,
		Skill:     skill,
		Modifiers: oneroll.Modifiers,
		Counter:   []int{1, 2, 3, 4, 5, 6, 7, 8},
		Capacities: map[string]float32{
			"Mass":  25.0,
			"Range": 10.0,
			"Speed": 2.5,
			"Self":  0.0,
		},
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_hyperskill.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		hsName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		ed, _ := strconv.Atoi(req.FormValue("Expert"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		skill.HyperSkill = new(oneroll.HyperSkill)

		hs = skill.HyperSkill

		hs.Name = hsName

		hs.Dice = &oneroll.DiePool{
			Normal: nd,
			Hard:   hd,
			Wiggle: wd,
			Expert: ed,
		}

		hs.Qualities = []*oneroll.Quality{}

		hs.Effect = req.FormValue("Effect")

		for _, qLoop := range wc.Counter[:4] { // Quality Loop

			qType := req.FormValue(fmt.Sprintf("Q%d-Type", qLoop))

			if qType != "" {
				l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-Level", qLoop)))
				if err != nil {
					l = 0
				}
				q := &oneroll.Quality{
					Type:  req.FormValue(fmt.Sprintf("Q%d-Type", qLoop)),
					Level: l,
					Name:  req.FormValue(fmt.Sprintf("Q%d-Name", qLoop)),
				}

				for _, cLoop := range wc.Counter[:4] {
					cType := req.FormValue(fmt.Sprintf("Q%d-C%d-Type", qLoop, cLoop))
					if cType != "" {
						cap := &oneroll.Capacity{
							Type: cType,
						}
						q.Capacities = append(q.Capacities, cap)
					}
				}

				q.Modifiers = []*oneroll.Modifier{}

				m := new(oneroll.Modifier)

				for _, mLoop := range wc.Counter { // Modifier Loop
					mName := req.FormValue(fmt.Sprintf("Q%d-M%d-Name", qLoop, mLoop))
					if mName != "" {
						l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-M%d-Level", qLoop, mLoop)))
						if err != nil {
							l = 0
						}

						m = oneroll.Modifiers[mName]

						if m.RequiresLevel {
							m.Level = l
						}

						if m.RequiresInfo {
							m.Info = req.FormValue(fmt.Sprintf("Q%d-M%d-Info", qLoop, mLoop))
						}
						q.Modifiers = append(q.Modifiers, m)
					}
				}
				hs.Qualities = append(hs.Qualities, q)
			}
		}

		apply := req.FormValue("Apply")

		if apply == "Yes" {
			// Apply Modifiers to Base
			hs.Apply = true

			for _, q := range hs.Qualities {
				for _, m := range q.Modifiers {

					skill.Modifiers = append(skill.Modifiers, m)
				}
			}
			hs.Effect += "\n++Added modifiers to base skill"
		}

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

// ModifyHyperSkillHandler renders a character in a Web page
func ModifyHyperSkillHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pk := vars["id"]
	sk := vars["skill"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Assign basic HyperStat
	skill := c.Skills[sk]

	hs := skill.HyperSkill

	// Assign additional empty Qualities to populate form
	if len(hs.Qualities) < 4 {
		for i := len(hs.Qualities); i < 4; i++ {
			tempQ := oneroll.NewQuality("")
			hs.Qualities = append(hs.Qualities, tempQ)
		}
	}

	// Assign additional empty Capacities to populate form
	for _, q := range hs.Qualities {
		if len(q.Capacities) < 4 {
			for i := len(q.Capacities); i < 4; i++ {
				tempC := oneroll.Capacity{
					Type: "",
				}
				q.Capacities = append(q.Capacities, &tempC)
			}
		}
		if len(q.Modifiers) < 8 {
			for i := len(q.Modifiers); i < 8; i++ {
				tempM := oneroll.NewModifier("")
				q.Modifiers = append(q.Modifiers, tempM)
			}
		}
	}

	wc := WebChar{
		Character: c,
		Skill:     skill,
		Modifiers: oneroll.Modifiers,
		Counter:   []int{1, 2, 3, 4, 5, 6, 7, 8},
		Capacities: map[string]float32{
			"Mass":  25.0,
			"Range": 10.0,
			"Speed": 2.5,
			"Self":  0.0,
		},
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_hyperskill.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		hsName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		ed, _ := strconv.Atoi(req.FormValue("Expert"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		skill.HyperSkill = new(oneroll.HyperSkill)

		hs = skill.HyperSkill

		hs.Name = hsName

		hs.Dice = &oneroll.DiePool{
			Normal: nd,
			Hard:   hd,
			Expert: ed,
			Wiggle: wd,
		}

		hs.Qualities = []*oneroll.Quality{}

		hs.Effect = req.FormValue("Effect")

		for _, qLoop := range wc.Counter[:4] { // Quality Loop

			qType := req.FormValue(fmt.Sprintf("Q%d-Type", qLoop))

			if qType != "" {
				l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-Level", qLoop)))
				if err != nil {
					l = 0
				}
				q := &oneroll.Quality{
					Type:  req.FormValue(fmt.Sprintf("Q%d-Type", qLoop)),
					Level: l,
					Name:  req.FormValue(fmt.Sprintf("Q%d-Name", qLoop)),
				}

				for _, cLoop := range wc.Counter[:4] {
					cType := req.FormValue(fmt.Sprintf("Q%d-C%d-Type", qLoop, cLoop))
					if cType != "" {
						cap := &oneroll.Capacity{
							Type: cType,
						}
						q.Capacities = append(q.Capacities, cap)
					}
				}

				q.Modifiers = []*oneroll.Modifier{}

				m := new(oneroll.Modifier)

				for _, mLoop := range wc.Counter { // Modifier Loop
					mName := req.FormValue(fmt.Sprintf("Q%d-M%d-Name", qLoop, mLoop))
					if mName != "" {
						l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-M%d-Level", qLoop, mLoop)))
						if err != nil {
							l = 0
						}

						m = oneroll.Modifiers[mName]

						if m.RequiresLevel {
							m.Level = l
						}

						if m.RequiresInfo {
							m.Info = req.FormValue(fmt.Sprintf("Q%d-M%d-Info", qLoop, mLoop))
						}
						q.Modifiers = append(q.Modifiers, m)
					}
				}
				hs.Qualities = append(hs.Qualities, q)
			}
		}

		apply := req.FormValue("Apply")

		if apply == "Yes" {
			// Apply Modifiers to Base
			hs.Apply = true

			for _, q := range hs.Qualities {
				for _, m := range q.Modifiers {

					skill.Modifiers = append(skill.Modifiers, m)
				}
			}
			hs.Effect += "\n++Added modifiers to base skill"
		}

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

// DeleteHyperSkillHandler renders a character in a Web page
func DeleteHyperSkillHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pk := vars["id"]
	sk := vars["skill"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	targetSkill := c.Skills[sk]

	wc := WebChar{
		Character: c,
		Skill:     targetSkill,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_hyperskill.html", wc)

	}

	if req.Method == "POST" {

		targetSkill.HyperSkill = nil

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

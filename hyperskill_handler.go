package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// AddHyperSkillHandler renders a character in a Web page
func AddHyperSkillHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]
	s := vars["skill"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	// Load CharacterModel
	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	c := cm.Character

	// Assign basic HyperSkill
	skill := c.Skills[s]

	skill.HyperSkill = &oneroll.HyperSkill{
		Name: fmt.Sprintf("Hyper-%s", skill.Name),
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
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		Skill:          skill,
		Modifiers:      oneroll.Modifiers,
		Counter:        []int{1, 2, 3, 4, 5, 6, 7, 8},
		Capacities: map[string]float32{
			"Mass":  25.0,
			"Range": 10.0,
			"Speed": 2.5,
			"Self":  0.0,
			"Touch": 0.0,
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

						tM := oneroll.Modifiers[mName]

						m = &tM

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
			// Remove existing Skill Modifiers
			skill.Modifiers = nil
			// Get base description text
			baseEffectText := strings.Split(hs.Effect, "++")

			// Add modifiers to base Skill
			for _, q := range hs.Qualities {
				for _, m := range q.Modifiers {

					skill.Modifiers = append(skill.Modifiers, m)
				}
			}
			// Update Skill Cost
			oneroll.UpdateCost(skill)
			// Determine the difference from the base skill cost
			modSkillCost := skill.Cost - (skill.Dice.Normal * skill.CostPerDie)

			// If difference is positive, add to descriptive text
			if modSkillCost > 0 {
				newModText := fmt.Sprintf("\n++Added modifiers to base skill (%dpts)",
					modSkillCost)
				hs.Effect = baseEffectText[0] + newModText
			} else {
				hs.Effect = baseEffectText[0]
			}
		}

		fmt.Println(c)

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)

		url := fmt.Sprintf("/view_character/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// ModifyHyperSkillHandler renders a character in a Web page
func ModifyHyperSkillHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]
	s := vars["skill"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	// Load CharacterModel
	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	c := cm.Character

	// Assign basic HyperStat
	skill := c.Skills[s]

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
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		Skill:          skill,
		Modifiers:      oneroll.Modifiers,
		Counter:        []int{1, 2, 3, 4, 5, 6, 7, 8},
		Capacities: map[string]float32{
			"Mass":  25.0,
			"Range": 10.0,
			"Speed": 2.5,
			"Self":  0.0,
			"Touch": 0.0,
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

						tM := oneroll.Modifiers[mName]

						m = &tM

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
			// Remove existing Skill Modifiers
			skill.Modifiers = nil
			// Get base description text
			baseEffectText := strings.Split(hs.Effect, "++")

			// Add modifiers to base Skill
			for _, q := range hs.Qualities {
				for _, m := range q.Modifiers {

					skill.Modifiers = append(skill.Modifiers, m)
				}
			}
			// Update Skill Cost
			oneroll.UpdateCost(skill)
			// Determine the difference from the base skill cost
			modSkillCost := skill.Cost - (skill.Dice.Normal * skill.CostPerDie)

			// If difference is positive, add to descriptive text
			if modSkillCost > 0 {
				newModText := fmt.Sprintf("\n++Added modifiers to base skill (%dpts)",
					modSkillCost)
				hs.Effect = baseEffectText[0] + newModText
			} else {
				hs.Effect = baseEffectText[0]
			}
		}

		fmt.Println(c)

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)

		url := fmt.Sprintf("/view_character/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// DeleteHyperSkillHandler renders a character in a Web page
func DeleteHyperSkillHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	sessionMap := getUserSessionValues(session)

	username := sessionMap["username"]
	loggedIn := sessionMap["loggedin"]
	isAdmin := sessionMap["isAdmin"]

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]
	s := vars["skill"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	// Load CharacterModel
	cm, err := database.PKLoadCharacterModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == cm.Author.UserName {
		IsAuthor = true
	} else {
		http.Redirect(w, req, "/", 302)
	}

	c := cm.Character

	targetSkill := c.Skills[s]

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		SessionUser:    username,
		IsLoggedIn:     loggedIn,
		IsAdmin:        isAdmin,
		Skill:          targetSkill,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_hyperskill.html", wc)

	}

	if req.Method == "POST" {

		targetSkill.HyperSkill = nil

		fmt.Println(c)

		err = database.UpdateCharacterModel(db, cm)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)

		url := fmt.Sprintf("/view_character/%d", cm.ID)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

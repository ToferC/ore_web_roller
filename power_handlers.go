package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/sessions"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
	"github.com/toferc/ore_web_roller/models"
)

func PowerIndexHandler(w http.ResponseWriter, req *http.Request) {

	pows, err := database.ListPowerModels(db)
	if err != nil {
		panic(err)
	}

	Render(w, "templates/index_powers.html", pows)
}

// PowerHandler renders a character in a Web page
func PowerHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	username := ""

	// Get session User
	u := session.Values["username"]

	// Type assertation
	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		username = user
	}

	vars := mux.Vars(req)
	pk := vars["id"]

	if len(pk) == 0 {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	pm, err := database.PKLoadPowerModel(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Validate that User == Author
	IsAuthor := false

	if username == pm.Author.UserName {
		IsAuthor = true
	}

	wc := WebChar{
		PowerModel: pm,
		IsAuthor:   IsAuthor,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/view_power.html", wc)

	}

	if req.Method == "POST" {

		// Parse Form and redirect
	}
}

// AddPowerHandler renders a character in a Web page
func AddPowerHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	username := ""

	// Get session User
	u := session.Values["username"]

	// Type assertation
	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		username = user
	}

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]

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

	// Create default Power to populate page
	defaultPower := &oneroll.Power{
		Name: "",
		Dice: &oneroll.DiePool{
			Normal: 0,
			Hard:   0,
			Wiggle: 0,
		},
		Effect:    "",
		Qualities: []*oneroll.Quality{},
		Slug:      "",
	}

	// Map default Power to Character.Powers
	if len(c.Powers) == 0 {
		c.Powers = map[string]*oneroll.Power{"default": defaultPower}
	} else {
		c.Powers["default"] = defaultPower
	}

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		Modifiers:      oneroll.Modifiers,
		Counter:        []int{1, 2, 3, 4, 5, 6, 7, 8},
		Capacities: map[string]float32{
			"Mass":  25.0,
			"Range": 10.0,
			"Speed": 2.5,
			"Self":  0.0,
		},
		Power: c.Powers["default"],
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_power.html", wc)

	}

	if req.Method == "POST" { // POST

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
			Effect:    req.FormValue("Effect"),
			Qualities: []*oneroll.Quality{},
			Slug:      oneroll.ToSnakeCase(pName),
		}

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

				for _, cLoop := range wc.Counter[:3] {
					cType := req.FormValue(fmt.Sprintf("Q%d-C%d-Type", qLoop, cLoop))
					if cType != "" {
						cap := &oneroll.Capacity{
							Type: cType,
						}
						q.Capacities = append(q.Capacities, cap)
					}
				}

				m := new(oneroll.Modifier)

				for _, mLoop := range wc.Counter { // Modifier Loop
					mName := req.FormValue(fmt.Sprintf("Q%d-M%d-Name", qLoop, mLoop))
					if mName != "" {

						tM := oneroll.Modifiers[mName]

						m = &tM

						if m.RequiresLevel {
							l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-M%d-Level", qLoop, mLoop)))
							if err != nil {
								l = 1
							}
							m.Level = l
						}

						if m.RequiresInfo {
							m.Info = req.FormValue(fmt.Sprintf("Q%d-M%d-Info", qLoop, mLoop))
						}
						q.Modifiers = append(q.Modifiers, m)
					}
				}
				p.Qualities = append(p.Qualities, q)
			}
		}

		c.Powers[p.Slug] = &p
		delete(c.Powers, "default")

		// Insert power into App archive if user authorizes
		if req.FormValue("Archive") != "" {
			p.DeterminePowerCapacities()
			p.CalculateCost()

			pm := models.PowerModel{
				Power:  &p,
				Author: cm.Author,
			}
			database.SavePowerModel(db, &pm)
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

		http.Redirect(w, req, url, http.StatusFound)
	}
}

// ModifyPowerHandler renders a character in a Web page
func ModifyPowerHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	username := ""

	// Get session User
	u := session.Values["username"]

	// Type assertation
	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		username = user
	}

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]
	pow := vars["power"]

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

	// Assign existing Power, Qualities, Capacities & Modifiers
	p := c.Powers[pow]

	// Assign additional empty Qualities to populate form
	if len(p.Qualities) < 4 {
		for i := len(p.Qualities); i < 4; i++ {
			tempQ := oneroll.NewQuality("")
			p.Qualities = append(p.Qualities, tempQ)
		}
	} else {
		// Always create at least 2 Qualities
		for i := 0; i < 2; i++ {
			tempQ := oneroll.NewQuality("")
			p.Qualities = append(p.Qualities, tempQ)
		}
	}

	// Assign additional empty Capacities to populate form
	for _, q := range p.Qualities {
		if len(q.Capacities) < 3 {
			for i := len(q.Capacities); i < 3; i++ {
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
		Modifiers:      oneroll.Modifiers,
		Counter:        []int{1, 2, 3, 4, 5},
		Capacities: map[string]float32{
			"Mass":  25.0,
			"Range": 10.0,
			"Speed": 2.5,
			"Self":  0.0,
		},
		Power: p,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/modify_power.html", wc)

	}

	if req.Method == "POST" { // POST

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
			Effect:    req.FormValue("Effect"),
			Qualities: []*oneroll.Quality{},
			Slug:      oneroll.ToSnakeCase(pName),
		}

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

				for _, cLoop := range wc.Counter[:3] {
					cType := req.FormValue(fmt.Sprintf("Q%d-C%d-Type", qLoop, cLoop))
					if cType != "" {
						cap := &oneroll.Capacity{
							Type: cType,
						}
						q.Capacities = append(q.Capacities, cap)
					}
				}

				m := new(oneroll.Modifier)

				for _, mLoop := range wc.Counter { // Modifier Loop
					mName := req.FormValue(fmt.Sprintf("Q%d-M%d-Name", qLoop, mLoop))
					if mName != "" {

						// Take base modifier struct from Modifiers
						tM := oneroll.Modifiers[mName]

						m = &tM

						if m.RequiresLevel {
							// Ensure level is a number or set to 1
							l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-M%d-Level", qLoop, mLoop)))
							if err != nil {
								l = 1
							}
							m.Level = l
						}

						if m.RequiresInfo {
							m.Info = req.FormValue(fmt.Sprintf("Q%d-M%d-Info", qLoop, mLoop))
						}
						// Append new modifier to Quality Modifiers
						q.Modifiers = append(q.Modifiers, m)

					}
				}
				// Append Quality to Power Qualities
				p.Qualities = append(p.Qualities, q)
			}
		}

		// Remove the base Power if necessary
		if pName != pow {
			// Remove previous map of the power if the name has changed
			delete(c.Powers, pow)
		}

		// Add Power to Character Powers map
		c.Powers[p.Slug] = &p

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

// DeletePowerHandler renders a character in a Web page
func DeletePowerHandler(w http.ResponseWriter, req *http.Request) {

	// Get session values or redirect to Login
	session, err := sessions.Store.Get(req, "session")

	if err != nil {
		log.Println("error identifying session")
		http.Redirect(w, req, "/login/", 302)
		return
		// in case of error
	}

	// Prep for user authentication
	username := ""

	// Get session User
	u := session.Values["username"]

	// Type assertation
	if user, ok := u.(string); !ok {
	} else {
		fmt.Println(user)
		username = user
	}

	// Get variables from URL
	vars := mux.Vars(req)
	pk := vars["id"]
	pow := vars["power"]

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

	// Assign existing Power, Qualities, Capacities & Modifiers
	p := c.Powers[pow]

	wc := WebChar{
		CharacterModel: cm,
		IsAuthor:       IsAuthor,
		Power:          p,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_power.html", wc)

	}

	if req.Method == "POST" {

		delete(c.Powers, pow)

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

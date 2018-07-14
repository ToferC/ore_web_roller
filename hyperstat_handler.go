package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// AddHyperStatHandler renders a character in a Web page
func AddHyperStatHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pk := vars["id"]
	s := vars["stat"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Assign basic HyperStat
	stat := c.Statistics[s]

	stat.HyperStat = &oneroll.HyperStat{
		Dice: &oneroll.DiePool{
			Normal: 0,
			Hard:   0,
			Wiggle: 0,
		},
		Effect: "",
	}

	hs := stat.HyperStat

	hs.Qualities = []*oneroll.Quality{}

	qualities := []string{"Attack", "Defend", "Useful", ""}

	// HyperStats start with all qualities, +1 extra
	for _, qs := range qualities {
		q := &oneroll.Quality{
			Type:       qs,
			Level:      0,
			CostPerDie: 0,
		}

		// Add the completed quality to Power
		hs.Qualities = append(hs.Qualities, q)
	}

	// Add Capacities (Self) to all basic Qualities

	for _, q := range hs.Qualities {
		if q.Type != "" {
			cap := oneroll.Capacity{Type: "Self"}
			q.Capacities = append(q.Capacities, &cap)
		}
	}

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
		Statistic: stat,
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
		Render(w, "templates/add_hyperstat.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		hsName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		stat.HyperStat = new(oneroll.HyperStat)

		hs = stat.HyperStat

		hs.Name = hsName

		hs.Dice = &oneroll.DiePool{
			Normal: nd,
			Hard:   hd,
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

					stat.Modifiers = append(stat.Modifiers, m)
				}
			}
			hs.Effect += "\n++Added modifiers to base stat"
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

// ModifyHyperStatHandler renders a character in a Web page
func ModifyHyperStatHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pk := vars["id"]
	s := vars["stat"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	// Assign basic HyperStat
	stat := c.Statistics[s]

	hs := stat.HyperStat

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
		Statistic: stat,
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
		Render(w, "templates/modify_hyperstat.html", wc)

	}

	if req.Method == "POST" { // POST

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		hsName := req.FormValue("Name")

		nd, _ := strconv.Atoi(req.FormValue("Normal"))
		hd, _ := strconv.Atoi(req.FormValue("Hard"))
		wd, _ := strconv.Atoi(req.FormValue("Wiggle"))

		stat.HyperStat = new(oneroll.HyperStat)

		hs = stat.HyperStat

		hs.Name = hsName

		hs.Dice = &oneroll.DiePool{
			Normal: nd,
			Hard:   hd,
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

					stat.Modifiers = append(stat.Modifiers, m)
				}
			}
			hs.Effect += "\n++Added modifiers to base stat"
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

// DeleteHyperStatHandler renders a character in a Web page
func DeleteHyperStatHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pk := vars["id"]
	s := vars["stat"]

	id, err := strconv.Atoi(pk)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	targetStat := c.Statistics[s]

	wc := WebChar{
		Character: c,
		Statistic: targetStat,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_hyperstat.html", wc)

	}

	if req.Method == "POST" {

		targetStat.HyperStat = nil

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

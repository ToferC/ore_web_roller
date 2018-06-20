package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

type WebChar struct {
	Character   *oneroll.Character
	Modifiers   map[string]*oneroll.Modifier
	Sources     map[string]*oneroll.Source
	Permissions map[string]*oneroll.Permission
	Intrinsics  map[string]*oneroll.Intrinsic
	Capacities  map[string]float32
	Counter     []int
}

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

	wc := WebChar{
		Character: c,
		Modifiers: oneroll.Modifiers,
		Counter:   []int{1, 2, 3, 4},
		Capacities: map[string]float32{
			"Mass":  25.0,
			"Range": 10.0,
			"Speed": 2.5,
			"Self":  0.0,
		},
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/form_power.html", wc)

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
			Effect:    req.FormValue("Effect"),
			Qualities: []*oneroll.Quality{},
		}

		for _, qLoop := range wc.Counter { // Quality Loop

			qType := req.FormValue(fmt.Sprintf("Q%d-Type", qLoop))

			if qType != "" {
				l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-Level", qLoop)))
				if err != nil {
					l = 1
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
						l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%d-M%d-Level", qLoop, mLoop)))
						if err != nil {
							l = 1
						}

						m = oneroll.Modifiers[req.FormValue(fmt.Sprintf("Q%d-M%d-Name", qLoop, mLoop))]

						if m.RequiresLevel {
							m.Level = l
						}

						if m.RequiresInfo {
							m.Info = req.FormValue(fmt.Sprintf("Q%d-M%d-Name", qLoop, mLoop))
						}
						q.Modifiers = append(q.Modifiers, m)
					}
				}
				p.Qualities = append(p.Qualities, q)
			}
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

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

type WebChar struct {
	Character   *oneroll.Character
	Power       *oneroll.Power
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
		Character: c,
		Modifiers: oneroll.Modifiers,
		Counter:   []int{1, 2, 3, 4, 5, 6, 7, 8},
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
			Slug:      oneroll.ToSnakeCase(pName),
		}

		for _, qLoop := range wc.Counter[:4] { // Quality Loop

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
				p.Qualities = append(p.Qualities, q)
			}
		}

		c.Powers[p.Slug] = &p
		delete(c.Powers, "default")

		fmt.Println(c)

		err = database.UpdateCharacter(db, c)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Saved")
		}

		fmt.Println(c)
		http.Redirect(w, req, "/add_power/"+string(c.ID), http.StatusSeeOther)
	}
}

// ModifyPowerHandler renders a character in a Web page
func ModifyPowerHandler(w http.ResponseWriter, req *http.Request) {

	s := req.URL.Path[len("/modify_power/"):]

	sSlice := strings.Split(s, "/")

	ch, pow := sSlice[0], sSlice[1]

	id, err := strconv.Atoi(ch)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
	}

	p := c.Powers[pow]

	wc := WebChar{
		Character: c,
		Modifiers: oneroll.Modifiers,
		Counter:   []int{1, 2, 3, 4, 5},
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

	} else { // POST

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

		for _, qLoop := range p.Qualities { // Quality Loop

			qType := req.FormValue(fmt.Sprintf("Q%s-Type", qLoop))

			if qType != "" {
				l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%s-Level", qLoop)))
				if err != nil {
					l = 1
				}
				q := &oneroll.Quality{
					Type:  req.FormValue(fmt.Sprintf("Q%s-Type", qLoop)),
					Level: l,
					Name:  req.FormValue(fmt.Sprintf("Q%s-Name", qLoop)),
				}

				for _, cLoop := range qLoop.Capacities {
					cType := req.FormValue(fmt.Sprintf("Q%s-C%s-Type", qLoop, cLoop))
					if cType != "" {
						cap := &oneroll.Capacity{
							Type: cType,
						}
						q.Capacities = append(q.Capacities, cap)
					}
				}

				m := new(oneroll.Modifier)

				for _, mLoop := range qLoop.Modifiers { // Modifier Loop
					mName := req.FormValue(fmt.Sprintf("Q%s-M%s-Name", qLoop, mLoop))
					if mName != "" {
						l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Q%s-M%s-Level", qLoop, mLoop)))
						if err != nil {
							l = 1
						}

						m = oneroll.Modifiers[mName]

						if m.RequiresLevel {
							m.Level = l
						}

						if m.RequiresInfo {
							m.Info = req.FormValue(fmt.Sprintf("Q%s-M%s-Info", qLoop, mLoop))
						}
						q.Modifiers = append(q.Modifiers, m)
					}
				}
				p.Qualities = append(p.Qualities, q)
			}
		}

		c.Powers[p.Slug] = &p

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

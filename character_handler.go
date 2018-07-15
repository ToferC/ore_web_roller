package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

func CharacterIndexHandler(w http.ResponseWriter, req *http.Request) {

	characters, err := database.ListCharacters(db)
	if err != nil {
		panic(err)
	}

	Render(w, "templates/index_characters.html", characters)
}

// CharacterHandler renders a character in a Web page
func CharacterHandler(w http.ResponseWriter, req *http.Request) {

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

	wc := WebChar{
		Character: c,
		Counter:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/view_character.html", wc)

	}

	if req.Method == "POST" {

		// Parse Form and redirect
		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		if c.Setting != "RE" {
			wp, err := strconv.Atoi(req.FormValue("Willpower"))
			if err != nil {
				wp = c.Willpower
			}
			c.Willpower = wp
		}

		// Set up character for actual play
		c.InPlay = true

		for k, v := range c.HitLocations {
			for i := range v.Shock {
				v.Shock[i] = false
				if req.FormValue(fmt.Sprintf("%s-Shock-%d", k, i)) != "" {
					v.Shock[i] = true
				}
			}
			for i := range v.Kill {
				v.Kill[i] = false
				if req.FormValue(fmt.Sprintf("%s-Kill-%d", k, i)) != "" {
					v.Kill[i] = true
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
		// Render page
		Render(w, "templates/view_character.html", wc)
	}

}

// NewCharacterHandler renders a character in a Web page
func NewCharacterHandler(w http.ResponseWriter, req *http.Request) {

	c := &oneroll.Character{Setting: "WT"}

	vars := mux.Vars(req)
	setting := vars["setting"]

	switch setting {
	case "SR":
		c = oneroll.NewSRCharacter("")
	case "WT":
		c = oneroll.NewWTCharacter("")
	case "RE":
		c = oneroll.NewReignCharacter("")
	}

	if c.Setting != "RE" {
		a := c.Archetype

		// Assign additional empty Sources to populate form
		if len(a.Sources) < 4 {
			for i := len(a.Sources); i < 4; i++ {
				tempS := oneroll.Source{
					Type: "",
				}
				a.Sources = append(a.Sources, &tempS)
			}
		}

		// Assign additional empty Permissions to populate form
		if len(a.Permissions) < 4 {
			for i := len(a.Permissions); i < 4; i++ {
				tempP := oneroll.Permission{
					Type: "",
				}
				a.Permissions = append(a.Permissions, &tempP)
			}
		}

		// Assign additional empty Sources to populate form
		if len(a.Intrinsics) < 5 {
			for i := len(a.Intrinsics); i < 5; i++ {
				tempI := oneroll.Intrinsic{
					Name: "",
				}
				a.Intrinsics = append(a.Intrinsics, &tempI)
			}
		}

		// Assign additional empty HitLocations to populate form
		if len(c.HitLocations) < 10 {
			for i := len(c.HitLocations); i < 10; i++ {
				t := oneroll.Location{
					Name: "",
				}
				c.HitLocations["z"+string(i)] = &t
			}
		}

	}

	wc := WebChar{
		Character:   c,
		Modifiers:   oneroll.Modifiers,
		Counter:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Sources:     oneroll.Sources,
		Permissions: oneroll.Permissions,
		Intrinsics:  oneroll.Intrinsics,
		Advantages:  nil,
	}

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/add_character.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		c := &oneroll.Character{}

		setting := req.FormValue("Setting")

		switch setting {
		case "SR":
			c = oneroll.NewSRCharacter(req.FormValue("Name"))
		case "WT":
			c = oneroll.NewWTCharacter(req.FormValue("Name"))
		case "RE":
			c = oneroll.NewReignCharacter(req.FormValue("Name"))
		}

		if setting == "SR" || setting == "WT" {
			c.Archetype = &oneroll.Archetype{
				Type: req.FormValue("Archetype"),
			}
			for _, s := range wc.Counter { // Loop

				sType := req.FormValue(fmt.Sprintf("Source-%d", s))

				pType := req.FormValue(fmt.Sprintf("Permission-%d", s))

				iName := req.FormValue(fmt.Sprintf("Intrinsic-%d-Name", s))

				iInfo := req.FormValue(fmt.Sprintf("Intrinsic-%d-Info", s))

				if iName != "" {
					i := oneroll.Intrinsics[iName]
					l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Intrinsic-%d-Level", s)))
					if err != nil {
						l = 1
					}
					i.Level = l
					i.Info = iInfo
					c.Archetype.Intrinsics = append(c.Archetype.Intrinsics, &i)
				}

				if sType != "" {
					tS := oneroll.Sources[sType]
					c.Archetype.Sources = append(c.Archetype.Sources, &tS)
				}
				if pType != "" {
					tP := oneroll.Permissions[pType]
					c.Archetype.Permissions = append(c.Archetype.Permissions, &tP)
				}
			}
		}

		c.Description = req.FormValue("Description")

		for _, st := range c.StatMap {
			c.Statistics[st].Dice.Normal, _ = strconv.Atoi(req.FormValue(st))
		}

		for _, sk := range c.Skills {
			sk.Dice.Normal, _ = strconv.Atoi(req.FormValue(sk.Name))
			if sk.ReqSpec {
				sk.Specialization = req.FormValue(fmt.Sprintf("%s-Spec", sk.Name))
			}
		}

		// Hit locations - need to add new map or amend old one

		newHL := map[string]*oneroll.Location{}

		for i := range c.HitLocations {

			name := req.FormValue(fmt.Sprintf("%s-Name", i))

			if name != "" {

				boxes, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-Boxes", i)))
				lar, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-LAR", i)))
				har, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-HAR", i)))

				fmt.Println(name, boxes, lar, har)

				newHL[name] = &oneroll.Location{
					Name:   name,
					Boxes:  boxes,
					LAR:    lar,
					HAR:    har,
					HitLoc: []int{},
				}

				newHL[name].FillWounds()

				for j := 1; j < 11; j++ {
					if req.FormValue(fmt.Sprintf("%s-%d-loc", i, j)) != "" {
						newHL[name].HitLoc = append(newHL[name].HitLoc, j)
					}
				}
			}
		}

		fmt.Println(newHL)
		c.HitLocations = newHL

		err = database.SaveCharacter(db, c)
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

// ModifyCharacterHandler renders a character in a Web page
func ModifyCharacterHandler(w http.ResponseWriter, req *http.Request) {

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

	if c.Setting != "RE" {
		a := c.Archetype

		// Assign additional empty Sources to populate form
		if len(a.Sources) < 4 {
			for i := len(a.Sources); i < 4; i++ {
				tempS := oneroll.Source{
					Type: "",
				}
				a.Sources = append(a.Sources, &tempS)
			}
		}

		// Assign additional empty Permissions to populate form
		if len(a.Permissions) < 4 {
			for i := len(a.Permissions); i < 4; i++ {
				tempP := oneroll.Permission{
					Type: "",
				}
				a.Permissions = append(a.Permissions, &tempP)
			}
		}

		// Assign additional empty Sources to populate form
		if len(a.Intrinsics) < 5 {
			for i := len(a.Intrinsics); i < 5; i++ {
				tempI := oneroll.Intrinsic{
					Name: "",
				}
				a.Intrinsics = append(a.Intrinsics, &tempI)
			}
		}

		// Assign additional empty HitLocations to populate form
		if len(c.HitLocations) < 10 {
			for i := len(c.HitLocations); i < 10; i++ {
				t := oneroll.Location{
					Name: "",
				}
				c.HitLocations["z"+string(i)] = &t
			}
		}

	}

	wc := WebChar{
		Character:   c,
		Modifiers:   oneroll.Modifiers,
		Counter:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Sources:     oneroll.Sources,
		Permissions: oneroll.Permissions,
		Intrinsics:  oneroll.Intrinsics,
	}

	if req.Method == "GET" {

		// Render page

		Render(w, "templates/modify_character.html", wc)

	}

	if req.Method == "POST" {

		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		c.Name = req.FormValue("Name")

		if c.Setting != "RE" {

			c.Archetype = &oneroll.Archetype{
				Type: req.FormValue("Archetype"),
			}

			for _, s := range wc.Counter[:3] { // Loop

				sType := req.FormValue(fmt.Sprintf("Source-%d", s))

				pType := req.FormValue(fmt.Sprintf("Permission-%d", s))

				if sType != "" {
					tS := oneroll.Sources[sType]
					c.Archetype.Sources = append(c.Archetype.Sources, &tS)
				}
				if pType != "" {
					tP := oneroll.Permissions[pType]
					c.Archetype.Permissions = append(c.Archetype.Permissions, &tP)
				}
			}

			for _, s := range wc.Counter[:5] {
				iName := req.FormValue(fmt.Sprintf("Intrinsic-%d-Name", s))

				if iName != "" {
					i := oneroll.Intrinsics[iName]

					if i.RequiresLevel {
						l, err := strconv.Atoi(req.FormValue(fmt.Sprintf("Intrinsic-%d-Level", s)))
						if err != nil {
							l = 1
						}
						i.Level = l
					}

					if i.RequiresInfo {
						iInfo := req.FormValue(fmt.Sprintf("Intrinsic-%d-Info", s))
						i.Info = iInfo
					}

					c.Archetype.Intrinsics = append(c.Archetype.Intrinsics, &i)
				}
			}
		}

		c.Description = req.FormValue("Description")

		bw, err := strconv.Atoi(req.FormValue("BaseWill"))
		if err != nil {
			bw = c.BaseWill
		}

		c.BaseWill = bw

		wp, err := strconv.Atoi(req.FormValue("Willpower"))
		if err != nil {
			wp = c.Willpower
		}

		c.Willpower = wp

		for _, st := range c.StatMap {
			c.Statistics[st].Dice.Normal, _ = strconv.Atoi(req.FormValue(st))
		}

		for _, sk := range c.Skills {
			sk.Dice.Normal, _ = strconv.Atoi(req.FormValue(sk.Name))
			if sk.ReqSpec {
				sk.Specialization = req.FormValue(fmt.Sprintf("%s-Spec", sk.Name))
			}
		}

		// Hit locations - need to add new map or amend old one

		newHL := map[string]*oneroll.Location{}

		for i := range c.HitLocations {

			name := req.FormValue(fmt.Sprintf("%s-Name", i))

			if name != "" {
				boxes, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-Boxes", i)))
				lar, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-LAR", i)))
				har, _ := strconv.Atoi(req.FormValue(fmt.Sprintf("%s-HAR", i)))

				fmt.Println(name, boxes, lar, har)

				newHL[name] = &oneroll.Location{
					Name:   name,
					Boxes:  boxes,
					LAR:    lar,
					HAR:    har,
					HitLoc: []int{},
				}

				newHL[name].FillWounds()

				for j := 1; j < 11; j++ {
					if req.FormValue(fmt.Sprintf("%s-%d-loc", i, j)) != "" {
						newHL[name].HitLoc = append(newHL[name].HitLoc, j)
					}
				}
			}
		}

		fmt.Println(newHL)
		c.HitLocations = newHL

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

// DeleteCharacterHandler renders a character in a Web page
func DeleteCharacterHandler(w http.ResponseWriter, req *http.Request) {

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

	if req.Method == "GET" {

		// Render page
		Render(w, "templates/delete_character.html", c)

	}

	if req.Method == "POST" {

		database.DeleteCharacter(db, c.ID)

		fmt.Println("Deleted ", c.Name)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

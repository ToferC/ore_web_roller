package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

const baseDieString string = "1ac+4d+0hd+0wd+0gf+0sp+1nr"
const blankDieString string = "ac+d+hd+wd+gf+sp+nr"

// RollHandler generates a Web user interface
func RollHandler(w http.ResponseWriter, req *http.Request) {

	pk := mux.Vars(req)["id"]

	var dieString string

	dieString = fmt.Sprintf("%sac+%sd+%shd+%swd+%sgf+%ssp+%snr+%sed",
		req.FormValue("ac"),
		req.FormValue("nd"),
		req.FormValue("hd"),
		req.FormValue("wd"),
		req.FormValue("gf"),
		req.FormValue("sp"),
		req.FormValue("nr"),
		req.FormValue("ed"),
	)

	if dieString == blankDieString {
		dieString = baseDieString
	}

	id, err := strconv.Atoi(pk)
	if err != nil {
		id = 9999
	}

	c, err := database.PKLoadCharacter(db, int64(id))
	if err != nil {
		fmt.Println(err)
		c = &oneroll.Character{
			Name: "Player 1",
		}
	}

	roll := oneroll.Roll{
		Actor:  c,
		Action: "Act",
	}

	nd, hd, wd, ed, gf, sp, ac, nr, err := roll.ParseString(dieString)

	wv := WebView{
		Actor:       []*oneroll.Character{c},
		Rolls:       []oneroll.Roll{},
		Matches:     []oneroll.Match{},
		Normal:      []int{nd},
		Hard:        []int{hd},
		Wiggle:      []int{wd},
		Expert:      []int{ed},
		GoFirst:     []int{gf},
		Spray:       []int{sp},
		Actions:     []int{ac},
		NumRolls:    []int{nr},
		ErrorString: []error{err},
		DieString:   []string{dieString},
	}

	for x := 0; x < wv.NumRolls[0]; x++ {
		tempRoll := roll
		tempRoll.Resolve(dieString)
		wv.Rolls = append(wv.Rolls, tempRoll)
		tempRoll = oneroll.Roll{}
	}

	if req.Method == "GET" {

		Render(w, "templates/roller.html", wv)

		// wv.Rolls = []oneroll.Roll{}

	}

	if req.Method == "POST" {

		ndQ := req.FormValue("nd")
		hdQ := req.FormValue("hd")
		wdQ := req.FormValue("wd")
		edQ := req.FormValue("ed")

		gfQ := req.FormValue("gofirst")
		spQ := req.FormValue("spray")
		acQ := req.FormValue("actions")

		nrQ := req.FormValue("numrolls")

		url := fmt.Sprintf("/roll/%d?ac=%s&gf=%s&hd=%s&nd=%s&nr=%s&sp=%s&wd=%s&ed=%s",
			c.ID,
			acQ,
			gfQ, // Update roll mechanism to use Modifiers GF
			hdQ,
			ndQ,
			nrQ, // Update roll mechanism to use Modifiers NR
			spQ, // Update roll mechanism to use Modifiers SP
			wdQ,
			edQ,
		)

		http.Redirect(w, req, url, http.StatusSeeOther)
	}
}

// OpposeHandler generates a Web user interface
func OpposeHandler(w http.ResponseWriter, req *http.Request) {

	var nd, hd, wd, gf, sp, ac, action string
	var nd2, hd2, wd2, gf2, sp2, ac2, action2 string

	if req.Method == "GET" {

		charString, charString2 := req.FormValue("name"), req.FormValue("name2")

		var dieString, dieString2 string

		dieString = fmt.Sprintf("%sac+%sd+%shd+%swd+%sgf+%ssp+%snr",
			req.FormValue("ac"),
			req.FormValue("nd"),
			req.FormValue("hd"),
			req.FormValue("wd"),
			req.FormValue("gf"),
			req.FormValue("sp"),
			req.FormValue("nr"),
		)

		dieString2 = fmt.Sprintf("%sac+%sd+%shd+%swd+%sgf+%ssp+%snr",
			req.FormValue("ac2"),
			req.FormValue("nd2"),
			req.FormValue("hd2"),
			req.FormValue("wd2"),
			req.FormValue("gf2"),
			req.FormValue("sp2"),
			req.FormValue("nr2"),
		)

		if dieString == blankDieString {
			dieString = baseDieString
		}

		if charString == "" {
			charString = "Player1"
		}

		if dieString2 == blankDieString {
			dieString2 = baseDieString
		}

		if charString2 == "" {
			charString2 = "Player2"
		}

		c := oneroll.Character{
			ID:   int64(9998),
			Name: charString,
		}

		d := oneroll.Character{
			ID:   int64(9999),
			Name: charString2,
		}

		roll := oneroll.Roll{
			Actor:  &c,
			Action: "Act",
		}

		roll2 := oneroll.Roll{
			Actor:  &d,
			Action: "Act",
		}

		nd, hd, wd, ed, gf, sp, ac, nr, _ := roll.ParseString(dieString)

		nd2, hd2, wd2, ed2, gf2, sp2, ac2, nr2, _ := roll.ParseString(dieString2)

		wv := WebView{
			Actor:     []*oneroll.Character{&c, &d},
			Rolls:     []oneroll.Roll{},
			Matches:   []oneroll.Match{},
			Normal:    []int{nd, nd2},
			Hard:      []int{hd, hd2},
			Wiggle:    []int{wd, wd2},
			Expert:    []int{ed, ed2},
			GoFirst:   []int{gf, gf2},
			Spray:     []int{sp, sp2},
			Actions:   []int{ac, ac2},
			NumRolls:  []int{nr, nr2},
			DieString: []string{dieString, dieString2},
		}

		roll.Resolve(dieString)
		wv.Rolls = append(wv.Rolls, roll)

		roll2.Resolve(dieString2)
		wv.Rolls = append(wv.Rolls, roll2)

		// Figure this out - what is an opposed roll and
		// How do we pass to web view
		wv.Matches = oneroll.OpposedRoll(&roll, &roll2)

		Render(w, "templates/opposed.html", wv)

		// wv.Rolls = []oneroll.Roll{}

	}

	if req.Method == "POST" {
		// Submit form

		// Player 1
		c := oneroll.Character{
			Name: req.FormValue("name"),
		}

		action = req.FormValue("action")

		nd = req.FormValue("nd")
		hd = req.FormValue("hd")
		wd = req.FormValue("wd")

		gf = req.FormValue("gofirst")
		sp = req.FormValue("spray")
		ac = req.FormValue("actions")

		// Player 2
		d := oneroll.Character{
			Name: req.FormValue("name2"),
		}

		action2 = req.FormValue("action2") // returns on/off

		nd2 = req.FormValue("nd2")
		hd2 = req.FormValue("hd2")
		wd2 = req.FormValue("wd2")

		gf2 = req.FormValue("gofirst2")
		sp2 = req.FormValue("spray2")
		ac2 = req.FormValue("actions2")

		fmt.Println(action, action2)

		q := req.URL.Query()

		q.Add("name1", c.Name)
		q.Add("ac", ac)
		q.Add("nd", nd)
		q.Add("hd", hd)
		q.Add("wd", wd)
		q.Add("gf", gf)
		q.Add("sp", sp)

		q.Add("name2", d.Name)
		q.Add("ac2", ac2)
		q.Add("nd2", nd2)
		q.Add("hd2", hd2)
		q.Add("wd2", wd2)
		q.Add("gf2", gf2)
		q.Add("sp2", sp2)

		// Encode URL.Query
		qs := q.Encode()

		http.Redirect(w, req, "/opposed/?"+qs, http.StatusSeeOther)
	}
}

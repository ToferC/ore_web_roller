package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/toferc/oneroll"
)

// WebView is a ontainer for Web_gui data
type WebView struct {
	Rolls       []oneroll.Roll
	Matches     []oneroll.Match
	Actor       []string
	Normal      []int
	Hard        []int
	Wiggle      []int
	GoFirst     []int
	Spray       []int
	Actions     []int
	NumRolls    []int
	DieString   []string
	ErrorString []error
}

const baseDieString string = "1ac+4d+0hd+0wd+0gf+0sp+1nr"
const blankDieString string = "ac+d+hd+wd+gf+sp+nr"

// RollHandler generates a Web user interface
func RollHandler(w http.ResponseWriter, req *http.Request) {

	var nd, hd, wd, gf, sp, ac string

	if req.Method == "GET" {

		fullString := req.URL.Path[len("/roll/"):]

		q, err := url.ParseQuery(fullString)
		if err != nil {
			log.Fatal(err)
		}

		charString := q.Get("name")

		var dieString string

		dieString = fmt.Sprintf("%sac+%sd+%shd+%swd+%sgf+%ssp+%snr",
			q.Get("ac"),
			q.Get("nd"),
			q.Get("hd"),
			q.Get("wd"),
			q.Get("gf"),
			q.Get("sp"),
			q.Get("nr"),
		)

		if dieString == blankDieString {
			dieString = baseDieString
		}

		if charString == "" {
			charString = "Player"
		}

		c := oneroll.Character{
			Name: charString,
		}

		roll := oneroll.Roll{
			Actor:  &c,
			Action: "Act",
		}

		nd, hd, wd, gf, sp, ac, nr, err := roll.ParseString(dieString)

		wv := WebView{
			Actor:       []string{charString},
			Rolls:       []oneroll.Roll{},
			Matches:     []oneroll.Match{},
			Normal:      []int{nd},
			Hard:        []int{hd},
			Wiggle:      []int{wd},
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

		render(w, "templates/roller.html", wv)

		// wv.Rolls = []oneroll.Roll{}

	} else {

		c := oneroll.Character{
			Name: req.FormValue("name"),
		}

		nd = req.FormValue("nd")
		hd = req.FormValue("hd")
		wd = req.FormValue("wd")

		gf = req.FormValue("gofirst")
		sp = req.FormValue("spray")
		ac = req.FormValue("actions")

		nr := req.FormValue("numrolls")

		q := req.URL.Query()

		q.Add("name", c.Name)
		q.Add("ac", ac)
		q.Add("nd", nd)
		q.Add("hd", hd)
		q.Add("wd", wd)
		q.Add("gf", gf)
		q.Add("sp", sp)
		q.Add("nr", nr)

		qs := q.Encode()

		http.Redirect(w, req, "/roll/"+qs, http.StatusSeeOther)
	}
}

// OpposeHandler generates a Web user interface
func OpposeHandler(w http.ResponseWriter, req *http.Request) {

	var nd, hd, wd, gf, sp, ac, action string
	var nd2, hd2, wd2, gf2, sp2, ac2, action2 string

	if req.Method == "GET" {

		//fullString := req.URL.Path[len("/opposed/"):]

		fullString := req.URL.Path[len("/opposed/"):]

		q, err := url.ParseQuery(fullString)
		if err != nil {
			log.Fatal(err)
		}

		charString, charString2 := q.Get("name"), q.Get("name2")

		var dieString, dieString2 string

		dieString = fmt.Sprintf("%sac+%sd+%shd+%swd+%sgf+%ssp+%snr",
			q.Get("ac"),
			q.Get("nd"),
			q.Get("hd"),
			q.Get("wd"),
			q.Get("gf"),
			q.Get("sp"),
			q.Get("nr"),
		)

		dieString2 = fmt.Sprintf("%sac+%sd+%shd+%swd+%sgf+%ssp+%snr",
			q.Get("ac2"),
			q.Get("nd2"),
			q.Get("hd2"),
			q.Get("wd2"),
			q.Get("gf2"),
			q.Get("sp2"),
			q.Get("nr2"),
		)

		if dieString == blankDieString {
			dieString = baseDieString
		}

		if charString == "" {
			charString = "Player 1"
		}

		if dieString2 == blankDieString {
			dieString2 = baseDieString
		}

		if charString2 == "" {
			charString2 = "Player 2"
		}

		c := oneroll.Character{
			Name: charString,
		}

		d := oneroll.Character{
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

		nd, hd, wd, gf, sp, ac, nr, _ := roll.ParseString(dieString)

		nd2, hd2, wd2, gf2, sp2, ac2, nr2, _ := roll.ParseString(dieString2)

		wv := WebView{
			Actor:       []string{charString, charString2},
			Rolls:       []oneroll.Roll{},
			Matches:     []oneroll.Match{},
			Normal:      []int{nd, nd2},
			Hard:        []int{hd, hd2},
			Wiggle:      []int{wd, wd2},
			GoFirst:     []int{gf, gf2},
			Spray:       []int{sp, sp2},
			Actions:     []int{ac, ac2},
			NumRolls:    []int{nr, nr2},
			ErrorString: []error{err},
			DieString:   []string{dieString, dieString2},
		}

		roll.Resolve(dieString)
		wv.Rolls = append(wv.Rolls, roll)

		roll2.Resolve(dieString2)
		wv.Rolls = append(wv.Rolls, roll2)

		// Figure this out - what is an opposed roll and
		// How do we pass to web view
		wv.Matches = oneroll.OpposedRoll(&roll, &roll2)

		render(w, "templates/opposed.html", wv)

		// wv.Rolls = []oneroll.Roll{}

	} else {
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

		q.Add("name", c.Name)
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

		qs := q.Encode()

		http.Redirect(w, req, "/opposed/"+qs, http.StatusSeeOther)
	}
}

// SplitLines transfomrs results text string into slice
func SplitLines(s string) []string {
	sli := strings.Split(s, "/n")
	return sli
}

func render(w http.ResponseWriter, filename string, data interface{}) {

	tmpl := make(map[string]*template.Template)

	//tplFuncMap := make(template.FuncMap)

	//tplFuncMap["SplitLines"] = SplitLines

	baseTemplate := "templates/layout.html"

	tmpl[filename] = template.Must(template.ParseFiles(filename, baseTemplate))

	if err := tmpl[filename].ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

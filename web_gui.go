package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/toferc/oneroll"
)

// WebView is a ontainer for Web_gui data
type WebView struct {
	Rolls       []oneroll.Roll
	Actor       string
	Normal      int
	Hard        int
	Wiggle      int
	GoFirst     int
	Spray       int
	Actions     int
	NumRolls    int
	DieString   string
	ErrorString error
}

// RollHandler generates a Web user interface
func RollHandler(w http.ResponseWriter, req *http.Request) {

	var nd, hd, wd, gf, sp, ac string

	if req.Method == "GET" {

		dieString := req.URL.Path[len("/roll/"):]

		if dieString == "" {
			dieString = "1ac+4d+0hd+0wd+0gf+0sp+1nr/name=Player"
		}

		charString := strings.Split(dieString, "/name=")

		if len(charString) < 2 {
			charString = []string{"", "Player"}
		}

		c := oneroll.Character{
			Name: charString[1],
		}

		roll := oneroll.Roll{
			Actor:  &c,
			Action: "Act",
		}

		nd, hd, wd, gf, sp, ac, nr, err := roll.ParseString(dieString)

		if err != nil {

		}

		wv := WebView{
			Actor:       charString[1],
			Rolls:       []oneroll.Roll{},
			Normal:      nd,
			Hard:        hd,
			Wiggle:      wd,
			GoFirst:     gf,
			Spray:       sp,
			Actions:     ac,
			NumRolls:    nr,
			ErrorString: err,
			DieString:   dieString,
		}

		for x := 0; x < wv.NumRolls; x++ {
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

		var text string

		text = fmt.Sprintf("%sac+%sd+%shd+%swd+%sgf+%ssp+%snr",
			ac,
			nd,
			hd,
			wd,
			gf,
			sp,
			nr,
		)

		http.Redirect(w, req, "/roll/"+text+"/name="+c.Name, http.StatusSeeOther)
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

package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/toferc/oneroll"
)

// WebView is a container for Web_gui data
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

// SplitLines transfomrs results text string into slice
func SplitLines(s string) []string {
	sli := strings.Split(s, "/n")
	return sli
}

func Render(w http.ResponseWriter, filename string, data interface{}) {

	tmpl := make(map[string]*template.Template)

	//tplFuncMap := make(template.FuncMap)

	//tplFuncMap["SplitLines"] = SplitLines

	baseTemplate := "templates/layout.html"

	tmpl[filename] = template.Must(template.ParseFiles(filename, baseTemplate))

	if err := tmpl[filename].ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

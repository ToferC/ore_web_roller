package main

import (
	"fmt"
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

func skillRoll(c *oneroll.Character, s string, ac int) string {

	rollString := fmt.Sprintf("ac=%d&gf=%d&hd=%d&name=%s&nd=%d&nr=%d&sp=%d&wd=%d",
		ac,
		0, // Update roll mechanism to use Modifiers
		c.Skills[s].Dice.Hard+c.Skills[s].LinkStat.Dice.Hard,
		c.Name,
		c.Skills[s].Dice.Normal+c.Skills[s].LinkStat.Dice.Normal,
		0, // Update roll mechanism to use Modifiers
		0, // Update roll mechanism to use Modifiers
		c.Skills[s].Dice.Wiggle+c.Skills[s].LinkStat.Dice.Wiggle,
	)
	return "/roll/" + rollString
}

func statRoll(c *oneroll.Character, s string, ac int) string {

	rollString := fmt.Sprintf("ac=%d&gf=%d&hd=%d&name=%s&nd=%d&nr=%d&sp=%d&wd=%d",
		ac,
		0, // Update roll mechanism to use Modifiers
		c.Statistics[s].Dice.Hard,
		c.Name,
		c.Statistics[s].Dice.Normal,
		0, // Update roll mechanism to use Modifiers
		0, // Update roll mechanism to use Modifiers
		c.Statistics[s].Dice.Wiggle,
	)
	return "/roll/" + rollString
}

func qualityRoll(c *oneroll.Character, p *oneroll.Power, q *oneroll.Quality, ac int) string {

	rollString := fmt.Sprintf("ac=%d&gf=%d&hd=%d&name=%s&nd=%d&nr=%d&sp=%d&wd=%d",
		ac,
		0, // Update roll mechanism to use Modifiers GF
		p.Dice.Hard,
		c.Name,
		p.Dice.Normal,
		0, // Update roll mechanism to use Modifiers NR
		0, // Update roll mechanism to use Modifiers SP
		p.Dice.Wiggle,
	)
	return "/roll/" + rollString
}

func subtract(a, b int) int {
	return a - b
}

func add(a, b int) int {
	return a + b
}

func Render(w http.ResponseWriter, filename string, data interface{}) {

	tmpl := make(map[string]*template.Template)

	// Set up FuncMap
	funcMap := template.FuncMap{
		"skillRoll":   skillRoll,
		"statRoll":    statRoll,
		"qualityRoll": qualityRoll,
		"subtract":    subtract,
		"add":         add,
	}

	baseTemplate := "templates/layout.html"

	tmpl[filename] = template.Must(template.New("").Funcs(funcMap).ParseFiles(filename, baseTemplate))

	if err := tmpl[filename].ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

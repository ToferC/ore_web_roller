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
	Actor       []*oneroll.Character
	Normal      []int
	Hard        []int
	Wiggle      []int
	Expert      []int
	GoFirst     []int
	Spray       []int
	Actions     []int
	DieString   []string
	NumRolls    []int
	ErrorString []error
}

type WebLoc struct {
}

// WebChar is a framework to send objects & data to a Web view
type WebChar struct {
	Character   *oneroll.Character
	Power       *oneroll.Power
	Statistic   *oneroll.Statistic
	Skill       *oneroll.Skill
	Shock       map[string][]int
	Kill        map[string][]int
	Modifiers   map[string]oneroll.Modifier
	Sources     map[string]oneroll.Source
	Permissions map[string]oneroll.Permission
	Intrinsics  map[string]oneroll.Intrinsic
	Advantages  map[string]oneroll.Advantage
	Capacities  map[string]float32
	Powers      map[string]oneroll.Power
	Counter     []int
}

// SplitLines transfomrs results text string into slice
func SplitLines(s string) []string {
	sli := strings.Split(s, "/n")
	return sli
}

func skillRoll(c *oneroll.Character, sk *oneroll.Skill, st *oneroll.Statistic, ac int) string {

	skill := oneroll.ReturnDice(sk)
	stat := oneroll.ReturnDice(st)

	normal := stat.Normal + skill.Normal
	hard := stat.Hard + skill.Hard
	expert := skill.Expert
	wiggle := stat.Wiggle + skill.Wiggle
	goFirst := oneroll.Max(stat.GoFirst, skill.GoFirst)
	spray := oneroll.Max(stat.Spray, skill.Spray)

	url := fmt.Sprintf("/roll/%d?ac=%d&gf=%d&hd=%d&nd=%d&nr=%d&sp=%d&wd=%d&ed=%d",
		c.ID,
		ac,
		goFirst,
		hard,
		normal,
		1, // Update roll mechanism to use Modifiers
		spray,
		wiggle,
		expert,
	)
	return url
}

func statRoll(c *oneroll.Character, s *oneroll.Statistic, ac int) string {

	td := oneroll.ReturnDice(s)

	normal := td.Normal
	hard := td.Hard
	wiggle := td.Wiggle
	goFirst := td.GoFirst
	spray := td.Spray

	url := fmt.Sprintf("/roll/%d?ac=%d&gf=%d&hd=%d&nd=%d&nr=%d&sp=%d&wd=%d",
		c.ID,
		ac,
		goFirst,
		hard,
		normal,
		1, // Update roll mechanism to use Modifiers
		spray,
		wiggle,
	)
	return url
}

func qualityRoll(c *oneroll.Character, p *oneroll.Power, q *oneroll.Quality, ac int) string {

	for _, m := range q.Modifiers {
		if m.Name == "Spray" {
			q.Dice.Spray = m.Level
		}

		if m.Name == "Go First" {
			q.Dice.GoFirst = m.Level
		}
	}

	url := fmt.Sprintf("/roll/%d?ac=%d&gf=%d&hd=%d&nd=%d&nr=%d&sp=%d&wd=%d",
		c.ID,
		ac,
		q.Dice.GoFirst, // Update roll mechanism to use Modifiers GF
		p.Dice.Hard,
		p.Dice.Normal,
		0,            // Update roll mechanism to use Modifiers NR
		q.Dice.Spray, // Update roll mechanism to use Modifiers SP
		p.Dice.Wiggle,
	)
	return url
}

func subtract(a, b int) int {
	return a - b
}

func add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

func isIn(s []int, t int) bool {
	for _, n := range s {
		if n == t {
			return true
		}
	}
	return false
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
		"multiply":    multiply,
		"isIn":        isIn,
	}

	baseTemplate := "templates/layout.html"

	tmpl[filename] = template.Must(template.New("").Funcs(funcMap).ParseFiles(filename, baseTemplate))

	if err := tmpl[filename].ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

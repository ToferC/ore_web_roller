package oneroll

import (
	"fmt"
	"strconv"
	"strings"
)

// ApplyDamage takes a match, weapon and and applies damage to locations
func ApplyDamage(c *Character, m Match, w Weapon) {

	height := m.Height
	width := m.Width

	var shock int
	var kill int

	// Parse damage code for shock
	if w.Shock == "" {
		strings.ToLower(w.Shock)

		switch {
		case strings.HasPrefix(w.Shock, "w") && len(w.Shock) > 1:
			n := strings.Split(w.Shock, "+")
			num, _ := strconv.Atoi(n[1])
			shock = width + num
		case w.Shock == "w":
			shock = width
		default:
			n := strings.Split(w.Shock, "+")
			num, _ := strconv.Atoi(n[1])
			shock = num
		}
	}

	// Parse damage code for Killing
	if w.Kill == "" {
		strings.ToLower(w.Kill)

		switch {
		case strings.HasPrefix(w.Kill, "w") && len(w.Kill) > 1:
			n := strings.Split(w.Kill, "+")
			num, _ := strconv.Atoi(n[1])
			shock = width + num
		case w.Kill == "w":
			shock = width
		default:
			n := strings.Split(w.Kill, "+")
			num, _ := strconv.Atoi(n[1])
			shock = num
		}
	}

	// Find locations
	fmt.Println(height, shock, kill)
	// Apply damage

}

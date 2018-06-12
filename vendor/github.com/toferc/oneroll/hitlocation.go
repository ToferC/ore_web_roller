package oneroll

import "fmt"

// Location represents a body area that can take damage
type Location struct {
	Name     string
	HitLoc   []int
	Boxes    int
	Stun     int
	Kill     int
	LAR      int
	HAR      int
	Disabled bool
}

// Strings
func (l Location) String() string {
	text := fmt.Sprintf("%s - %s: Boxes: %d",
		TrimSliceBrackets(l.HitLoc),
		l.Name,
		l.Boxes,
	)

	if l.LAR > 0 {
		text += fmt.Sprintf(" LAR %d", l.LAR)
	}

	if l.HAR > 0 {
		text += fmt.Sprintf(" HAR %d", l.HAR)
	}

	if l.Kill > 0 {
		text += fmt.Sprintf(" Kill %d", l.Kill)
	}

	if l.Stun > 0 {
		text += fmt.Sprintf(" Stun %d", l.Stun)
	}
	return text
}

package oneroll

import "fmt"

type Advantage struct {
	Name         string
	Level        int
	CostPerLevel int
	Description  string
}

func (a *Advantage) String() string {
	text := fmt.Sprintf("%s (%d)",
		a.Name, a.Level*a.CostPerLevel)

	return text
}

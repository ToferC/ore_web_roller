package oneroll

import "fmt"

// Statistic represents common attributes possessed by every character
type Statistic struct {
	Name    string
	Dice    *DiePool
	Booster []*Booster
}

// HyperStat is a modified version of a regular Statistic
type HyperStat struct {
	Name       string
	Dice       *DiePool
	Capacities []*Capacity
	Extras     []*Extra
	Flaws      []*Flaw
	CostPerDie int
	Booster    []*Booster
}

func (s Statistic) String() string {
	text := fmt.Sprintf("%s: %dd",
		s.Name,
		s.Dice.Normal,
	)

	if s.Dice.Hard > 0 {
		text += fmt.Sprintf("+%dhd", s.Dice.Hard)
	}

	if s.Dice.Wiggle > 0 {
		text += fmt.Sprintf("+%dwd", s.Dice.Wiggle)
	}

	if s.Dice.GoFirst > 0 {
		text += fmt.Sprintf(" Go First %d", s.Dice.GoFirst)
	}

	if s.Dice.Spray > 0 {
		text += fmt.Sprintf(" Spray %d", s.Dice.Spray)
	}

	return text
}

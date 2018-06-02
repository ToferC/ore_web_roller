package oneroll

import "fmt"

// Skill represents specific training
type Skill struct {
	Name           string
	LinkStat       *Statistic
	Dice           *DiePool
	ReqSpec        bool
	Specialization string
}

// HyperSkill is a modified version of a regular Skill
type HyperSkill struct {
	Name       string
	LinkStat   *Statistic
	Dice       *DiePool
	Capacities []*Capacity
	Extras     []*Extra
	Flaws      []*Flaw
	CostPerDie int
}

func (s Skill) String() string {

	text := fmt.Sprintf("%s ",
		s.Name)

	if s.ReqSpec {
		text += fmt.Sprintf("[%s] ", s.Specialization)
	}

	text += fmt.Sprintf("(%s): %dd",
		s.LinkStat.Name,
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

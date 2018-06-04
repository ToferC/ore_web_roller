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

	text += fmt.Sprintf(": %dd",
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

// ShowSkills shows skills grouped under stats
// all bool determines if all skills are shown or just the ones with dice in them.
func ShowSkills(c *Character, allSkills bool) string {
	statistics := []*Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

	var text string

	for _, stat := range statistics {
		text += fmt.Sprintf("%s\n", stat)
		for _, skill := range c.Skills {
			if skill.LinkStat.Name == stat.Name {
				if allSkills {
					// We want all skills
					text += fmt.Sprintf("-- %s\n", skill)
				} else {
					// We only want rated skills
					if SkillRated(skill) {
						text += fmt.Sprintf("-- %s\n", skill)
					}
				}
			}
		}
	}
	return text
}

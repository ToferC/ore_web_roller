package oneroll

import (
	"fmt"
)

// Character represents a full character in the ORE game
type Character struct {
	ID           int64
	Name         string
	Body         *Statistic
	Coordination *Statistic
	Sense        *Statistic
	Mind         *Statistic
	Command      *Statistic
	Charm        *Statistic
	BaseWill     int
	Willpower    int
	Skills       map[string]*Skill
	Archetype    *Archetype
	HyperStats   map[string]*HyperStat
	HyperSkills  map[string]*HyperSkill
	Permissions  map[string]*Permission
	Powers       map[string]*Power
	HitLocations map[string]*Location
	PointCost    int
}

// Display character
func (c *Character) String() string {

	statistics := []*Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

	text := fmt.Sprintf("\n%s (%d pts)\n", c.Name, c.PointCost)

	if c.Archetype.Type != "" {
		text += fmt.Sprint(c.Archetype)
	}

	text += "\n\nStats:\n"

	text += ShowSkills(c, false)

	text += fmt.Sprintf("\nBase Will:%d\n", c.BaseWill)
	text += fmt.Sprintf("Willpower: %d\n", c.Willpower)

	text += fmt.Sprintf("\nHit Locations:\n")

	for _, loc := range c.HitLocations {
		text += fmt.Sprintf("%s\n", loc)
	}

	if len(c.Archetype.Sources) > 0 {
		text += fmt.Sprintf("\nPowers:\n")

		for _, s := range statistics {
			if s.HyperStat != nil {
				text += fmt.Sprintf("%s\n\n", s.HyperStat)
			}
		}

		for _, s := range c.Skills {
			if s.HyperSkill != nil {
				text += fmt.Sprintf("%s\n\n", s.HyperSkill)
			}
		}

		for _, p := range c.Powers {
			text += fmt.Sprintf("%s\n\n", p)
		}
	}

	return text
}

// CalculateCost updates the character and sums
// total costs of all character elements. Call this on each character update
func (c *Character) CalculateCost() {

	var cost int

	if len(c.Archetype.Sources) > 0 {
		UpdateCost(c.Archetype)
		cost += c.Archetype.Cost
	}

	statistics := []*Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

	for _, stat := range statistics {
		UpdateCost(stat)
		cost += stat.Cost

		if stat.HyperStat != nil {
			UpdateCost(stat.HyperStat)
			cost += stat.HyperStat.Cost
		}
	}

	for _, skill := range c.Skills {
		UpdateCost(skill)
		cost += skill.Cost

		if skill.HyperSkill != nil {
			UpdateCost(skill.HyperSkill)
			cost += skill.HyperSkill.Cost
		}
	}

	for _, power := range c.Powers {
		// Determine power capacities
		power.DeterminePowerCapacities()
		UpdateCost(power)
		cost += power.Cost
	}

	comTotal := c.Command.Dice.Normal + c.Command.Dice.Hard + c.Command.Dice.Wiggle
	charmTotal := c.Charm.Dice.Normal + c.Charm.Dice.Hard + c.Charm.Dice.Wiggle

	cost += 3*c.BaseWill - (comTotal + charmTotal)
	cost += c.Willpower - c.BaseWill

	c.PointCost = cost
}

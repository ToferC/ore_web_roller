package oneroll

import (
	"fmt"
)

// Character represents a full character in the ORE game
type Character struct {
	ID           int64
	Name         string
	Setting      string
	Statistics   map[string]*Statistic
	StatMap      []string
	BaseWill     int
	Willpower    int
	Skills       map[string]*Skill
	Archetype    *Archetype
	HyperStats   map[string]*HyperStat
	HyperSkills  map[string]*HyperSkill
	Permissions  map[string]*Permission
	Powers       map[string]*Power
	HitLocations map[string]*Location
	LocationMap  []string
	PointCost    int
	InPlay       bool
}

// Display character
func (c *Character) String() string {

	text := fmt.Sprintf("\n%s (%d pts)\n", c.Name, c.PointCost)

	if c.Archetype.Type != "" {
		text += fmt.Sprint(c.Archetype)
	}

	text += "\n\nStats:\n"

	text += ShowSkills(c, false)

	text += fmt.Sprintf("\nBase Will: %d\n", c.BaseWill)
	text += fmt.Sprintf("Willpower: %d\n", c.Willpower)

	text += fmt.Sprintf("\nHit Locations:\n")

	for _, loc := range c.LocationMap {
		text += fmt.Sprintf("%s\n", c.HitLocations[loc])
	}

	if len(c.Archetype.Sources) > 0 {
		text += fmt.Sprintf("\nPowers:\n")

		for _, stat := range c.StatMap {
			s := c.Statistics[stat]
			if s.HyperStat != nil {
				text += fmt.Sprintf("\n%s\n", s.HyperStat)
				if len(s.Modifiers) > 0 {
					text += fmt.Sprintf("+ added modifiers to main stat: ")
					for _, m := range s.Modifiers {
						text += fmt.Sprintf("%s (%d/die) ", m.Name, m.Cost)
					}
				}
				text += fmt.Sprint("\n")
			}
		}

		for _, s := range c.Skills {
			if s.HyperSkill != nil {
				text += fmt.Sprintf("\n%s\n", s.HyperSkill)
				if len(s.Modifiers) > 0 {
					text += fmt.Sprintf("+ added modifiers to main stat: ")
					for _, m := range s.Modifiers {
						text += fmt.Sprintf("%s (%d/die) ", m.Name, m.Cost)
					}
				}
				text += fmt.Sprint("\n")
			}
		}

		for _, p := range c.Powers {
			text += fmt.Sprintf("%s", p)

			for _, q := range p.Qualities {
				text += fmt.Sprintln(q)
			}

			if p.Effect != "" {
				text += fmt.Sprintf("Effect: %s", p.Effect)
			}
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

	for _, stat := range c.Statistics {
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

	// Update BaseWill automaticallly if Character isn't in play

	calcBaseWill := 0

	for _, stat := range c.Statistics {
		if stat.EffectsWill {
			calcBaseWill += SumDice(stat.Dice)
			if stat.HyperStat != nil {
				calcBaseWill += SumDice(stat.HyperStat.Dice)
			}
		}
	}

	if !c.InPlay {
		c.BaseWill = calcBaseWill
		c.Willpower = calcBaseWill
	} else {
		cost += 3*c.BaseWill - calcBaseWill
		cost += c.Willpower - c.BaseWill
	}

	c.PointCost = cost
}

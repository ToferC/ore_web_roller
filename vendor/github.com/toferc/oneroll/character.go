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
	Archtypes    map[string]*Archtype
	HyperStats   map[string]*HyperStat
	HyperSkills  map[string]*HyperSkill
	Permissions  map[string]*Permission
	Powers       map[string]*Power
	HitLocations map[string]*Location
	PointCost    int
}

// NewWTCharacter generates an ORE WT character
func NewWTCharacter(name string) *Character {

	c := Character{
		Name: name,
	}

	c.Body = &Statistic{
		Name: "Body",
		Dice: &DiePool{
			Normal:  2,
			Hard:    0,
			GoFirst: 0,
		},
	}

	c.Coordination = &Statistic{
		Name: "Coordination",
		Dice: &DiePool{
			Normal: 2,
		},
	}
	c.Sense = &Statistic{
		Name: "Sense",
		Dice: &DiePool{
			Normal: 2,
		},
	}
	c.Mind = &Statistic{
		Name: "Mind",
		Dice: &DiePool{
			Normal: 2,
		},
	}
	c.Command = &Statistic{
		Name: "Command",
		Dice: &DiePool{
			Normal: 2,
		},
	}
	c.Charm = &Statistic{
		Name: "Charm",
		Dice: &DiePool{
			Normal: 2,
		},
	}

	c.HitLocations = map[string]*Location{
		"Head": &Location{
			Name:     "Head",
			HitLoc:   []int{10},
			Boxes:    4,
			Stun:     0,
			Kill:     0,
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Body": &Location{
			Name:     "Body",
			HitLoc:   []int{7, 8, 9},
			Boxes:    10,
			Stun:     0,
			Kill:     0,
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Left Arm": &Location{
			Name:     "Left Arm",
			HitLoc:   []int{5, 6},
			Boxes:    6,
			Stun:     0,
			Kill:     0,
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Right Arm": &Location{
			Name:     "Right Arm",
			HitLoc:   []int{3, 4},
			Boxes:    6,
			Stun:     0,
			Kill:     0,
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Left Leg": &Location{
			Name:     "Left Leg",
			HitLoc:   []int{2},
			Boxes:    6,
			Stun:     0,
			Kill:     0,
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Right Leg": &Location{
			Name:     "Right Leg",
			HitLoc:   []int{1},
			Boxes:    6,
			Stun:     0,
			Kill:     0,
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
	}

	c.BaseWill = c.Command.Dice.Normal + c.Charm.Dice.Normal
	c.Willpower = c.BaseWill

	c.Skills = map[string]*Skill{
		// Body Skills
		"Athletics": &Skill{
			Name:     "Athletics",
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Block": &Skill{
			Name:     "Block",
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Brawling": &Skill{
			Name:     "Brawling",
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Endurance": &Skill{
			Name:     "Endurance",
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Melee Weapon": &Skill{
			Name:     "Melee Weapon",
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
			ReqSpec:        true,
			Specialization: "Sword",
		},
		// Coordination Skills
		"Dodge": &Skill{
			Name:     "Dodge",
			LinkStat: c.Coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		"Driving": &Skill{
			Name:     "Driving",
			LinkStat: c.Coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
			ReqSpec:        true,
			Specialization: "Ground",
		},
		"Ranged Weapon": &Skill{
			Name:     "Ranged Weapon",
			LinkStat: c.Coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
			ReqSpec:        true,
			Specialization: "Pistol",
		},
		"Stealth": &Skill{
			Name:     "Stealth",
			LinkStat: c.Coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		// Sense Skills
		"Empathy": &Skill{
			Name:     "Empathy",
			LinkStat: c.Sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Perception": &Skill{
			Name:     "Perception",
			LinkStat: c.Sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Scrutiny": &Skill{
			Name:     "Scrutiny",
			LinkStat: c.Sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Mind Skills
		"First Aid": &Skill{
			Name:     "First Aid",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Knowledge": &Skill{
			Name:     "Knowledge",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Alchemy",
		},
		"Languages": &Skill{
			Name:     "Languages",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Chinese",
		},
		"Medicine": &Skill{
			Name:     "Medicine",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Navigation": &Skill{
			Name:     "Navigation",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Research": &Skill{
			Name:     "Research",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Security Systems": &Skill{
			Name:     "Security Systems",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Streetwise": &Skill{
			Name:     "Streetwise",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Survival": &Skill{
			Name:     "Survival",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Tactics": &Skill{
			Name:     "Tactics",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Charm Skills
		"Lie": &Skill{
			Name:     "Lie",
			LinkStat: c.Charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Performance": &Skill{
			Name:     "Performance",
			LinkStat: c.Charm,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Standup",
		},
		"Persuasion": &Skill{
			Name:     "Persuasion",
			LinkStat: c.Charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Command Skills
		"Interrogation": &Skill{
			Name:     "Interrogation",
			LinkStat: c.Command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Intimidation": &Skill{
			Name:     "Intimidation",
			LinkStat: c.Command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Leadership": &Skill{
			Name:     "Leadership",
			LinkStat: c.Command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Stability": &Skill{
			Name:     "Stability",
			LinkStat: c.Command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
	}

	return &c
}

// Display character
func (c *Character) String() string {

	text := fmt.Sprintf("\n%s\n\nStats:\n", c.Name)

	text += ShowSkills(c, false)

	text += fmt.Sprintf("\nBase Will:%d\n", c.BaseWill)
	text += fmt.Sprintf("Willpower: %d\n", c.Willpower)

	text += fmt.Sprintf("\nHit Locations:\n")

	for _, loc := range c.HitLocations {
		text += fmt.Sprintf("%s\n", loc)
	}
	return text
}

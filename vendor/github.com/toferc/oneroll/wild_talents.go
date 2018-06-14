package oneroll

// NewWTCharacter generates an ORE WT character
func NewWTCharacter(name string) *Character {

	c := Character{
		Name: name,
	}

	c.Archetype = new(Archetype)

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
			Name: "Athletics",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Block": &Skill{
			Name: "Block",
			Quality: &Quality{
				Type:  "Defend",
				Level: 1,
			},
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Brawling": &Skill{
			Name: "Brawling",
			Quality: &Quality{
				Type:  "Attack",
				Level: 1,
			},
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Endurance": &Skill{
			Name: "Endurance",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Melee Weapon": &Skill{
			Name: "Melee Weapon",
			Quality: &Quality{
				Type:  "Attack",
				Level: 1,
			},
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
			Name: "Dodge",
			Quality: &Quality{
				Type:  "Defend",
				Level: 1,
			},
			LinkStat: c.Coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		"Driving": &Skill{
			Name: "Driving",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
			ReqSpec:        true,
			Specialization: "Ground",
		},
		"Ranged Weapon": &Skill{
			Name: "Ranged Weapon",
			Quality: &Quality{
				Type:  "Attack",
				Level: 1,
			},
			LinkStat: c.Coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
			ReqSpec:        true,
			Specialization: "Pistol",
		},
		"Stealth": &Skill{
			Name: "Stealth",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		// Sense Skills
		"Empathy": &Skill{
			Name: "Empathy",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Perception": &Skill{
			Name: "Perception",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Scrutiny": &Skill{
			Name: "Scrutiny",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Mind Skills
		"First Aid": &Skill{
			Name: "First Aid",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Knowledge": &Skill{
			Name: "Knowledge",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Alchemy",
		},
		"Languages": &Skill{
			Name: "Languages",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Chinese",
		},
		"Medicine": &Skill{
			Name: "Medicine",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Navigation": &Skill{
			Name: "Navigation",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Research": &Skill{
			Name: "Research",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Security Systems": &Skill{
			Name: "Security Systems",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Streetwise": &Skill{
			Name: "Streetwise",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Survival": &Skill{
			Name: "Survival",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Tactics": &Skill{
			Name: "Tactics",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Charm Skills
		"Lie": &Skill{
			Name: "Lie",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Performance": &Skill{
			Name: "Performance",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Charm,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Standup",
		},
		"Persuasion": &Skill{
			Name: "Persuasion",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Command Skills
		"Interrogation": &Skill{
			Name: "Interrogation",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Intimidation": &Skill{
			Name: "Intimidation",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Leadership": &Skill{
			Name: "Leadership",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Stability": &Skill{
			Name: "Stability",
			Quality: &Quality{
				Type:  "Useful",
				Level: 1,
			},
			LinkStat: c.Command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
	}

	c.Powers = map[string]*Power{}

	return &c
}

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

// Statistic represents common attributes possessed by every character
type Statistic struct {
	Name    string
	Dice    *DiePool
	Booster []*Booster
}

// Skill represents specific training
type Skill struct {
	Name     string
	LinkStat *Statistic
	Dice     *DiePool
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

// Archtype is a grouping of Sources, Permissions & Intrinsics that defines what powers a character can use
type Archtype struct {
	Sources     []*Source
	Permissions []*Permission
	Intrinsics  []*Intrinsic
}

// Source is a source of a Character's powers
type Source struct {
	Type        string
	Cost        int
	Description string
}

// Permission is the type of powers a Character can purchase
type Permission struct {
	Type        string
	Cost        int
	Description string
}

// Intrinsic is a modification from the human standard
type Intrinsic struct {
	Name        string
	Cost        int
	Description string
}

// Power is a non-standard ability or miracle
type Power struct {
	Name      string
	Qualities []*Quality
	Dice      *DiePool
	Effect    string
	Dud       bool
}

// Quality is either Attack, Defend or Useful
type Quality struct {
	Type       string
	Level      int
	Capacities []*Capacity
	Extras     []*Extra
	Flaws      []*Flaw
	CostPerDie int
}

// Capacity is Range, Mass, Touch or Speed
type Capacity struct {
	Type    string
	Level   int
	Booster *Booster
}

// Booster multiplies a Capacity or Statistic
type Booster struct {
	Level int
}

// Extra enhances the abilities of a Power Quality
type Extra struct {
	Name     string
	Modifier int
}

// Flaw limits the abilities of a Power Quality
type Flaw struct {
	Name     string
	Modifier int
}

// Location represents a body area that can take damage
type Location struct {
	Name         string
	HitLoc       []int
	Boxes        int
	CurrentStun  int
	CurrentWound int
	LAR          int
	HAR          int
	Disabled     bool
}

// NewCharacter generates an ORE character
func NewCharacter(name string) *Character {

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
			Name:         "Head",
			HitLoc:       []int{10},
			Boxes:        4,
			CurrentStun:  0,
			CurrentWound: 0,
			LAR:          0,
			HAR:          0,
			Disabled:     false,
		},
		"Body": &Location{
			Name:         "Body",
			HitLoc:       []int{7, 8, 9},
			Boxes:        10,
			CurrentStun:  0,
			CurrentWound: 0,
			LAR:          0,
			HAR:          0,
			Disabled:     false,
		},
		"Left Arm": &Location{
			Name:         "Left Arm",
			HitLoc:       []int{5, 6},
			Boxes:        6,
			CurrentStun:  0,
			CurrentWound: 0,
			LAR:          0,
			HAR:          0,
			Disabled:     false,
		},
		"Right Arm": &Location{
			Name:         "Right Arm",
			HitLoc:       []int{3, 4},
			Boxes:        6,
			CurrentStun:  0,
			CurrentWound: 0,
			LAR:          0,
			HAR:          0,
			Disabled:     false,
		},
		"Left Leg": &Location{
			Name:         "Left Leg",
			HitLoc:       []int{2},
			Boxes:        6,
			CurrentStun:  0,
			CurrentWound: 0,
			LAR:          0,
			HAR:          0,
			Disabled:     false,
		},
		"Right Leg": &Location{
			Name:         "Right Leg",
			HitLoc:       []int{1},
			Boxes:        6,
			CurrentStun:  0,
			CurrentWound: 0,
			LAR:          0,
			HAR:          0,
			Disabled:     false,
		},
	}

	c.BaseWill = c.Command.Dice.Normal + c.Charm.Dice.Normal
	c.Willpower = c.BaseWill

	c.Skills = map[string]*Skill{
		"Athletics": &Skill{
			Name:     "Athletics",
			LinkStat: c.Body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"SmallArms": &Skill{
			Name:     "SmallArms",
			LinkStat: c.Coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		"Perception": &Skill{
			Name:     "Perception",
			LinkStat: c.Sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Business": &Skill{
			Name:     "Business",
			LinkStat: c.Mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
	}

	return &c
}

// Display character
func (c *Character) Display() {

	fmt.Println(c.Name)

	stats := []*Statistic{c.Body, c.Coordination, c.Sense,
		c.Mind, c.Command, c.Charm}

	for _, stat := range stats {
		text := fmt.Sprintf("%dd+%dhd+%dwd+%dgf+%dsp",
			stat.Dice.Normal,
			stat.Dice.Hard,
			stat.Dice.Wiggle,
			stat.Dice.GoFirst,
			stat.Dice.Spray,
		)
		fmt.Printf("%s: %s\n", stat.Name, text)
	}
	for _, loc := range c.HitLocations {
		fmt.Println(loc)
	}
	for _, skill := range c.Skills {
		fmt.Println(skill.Name, skill.Dice.Normal)
		//fmt.Println(skill.Name, FormSkillDieString(skill, 1))
	}
}

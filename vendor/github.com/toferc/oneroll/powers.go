package oneroll

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

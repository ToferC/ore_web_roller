package main

import (
	"fmt"

	"github.com/toferc/oneroll"
)

// Sets interface on dieFormat() behavior
type formatter interface {
	dieFormat() string
}

type St struct {
	Dice *oneroll.DiePool
}

func (s *St) dieFormat() string {

	actions := 1
	normal := s.Dice.Normal
	hard := s.Dice.Hard
	wiggle := s.Dice.Wiggle
	goFirst := s.Dice.GoFirst
	spray := s.Dice.Spray

	text := fmt.Sprintf("%dac+%dd+%dhd+%dwd+%dgf+%dsp",
		actions,
		normal,
		hard,
		wiggle,
		goFirst,
		spray)

	return text
}

type Po struct {
	Dice *oneroll.DiePool
}

func (p *Po) dieFormat() string {

	actions := 1
	normal := p.Dice.Normal
	hard := p.Dice.Hard
	wiggle := p.Dice.Wiggle
	goFirst := p.Dice.GoFirst
	spray := p.Dice.Spray

	text := fmt.Sprintf("%dac+%dd+%dhd+%dwd+%dgf+%dsp",
		actions,
		normal,
		hard,
		wiggle,
		goFirst,
		spray)

	return text
}

// Call printFormat from St or Po to call dieFormat
func printFormat(f formatter) string {
	return f.dieFormat()
}

func main() {
	p := Po{Dice: &oneroll.DiePool{
		Normal: 4,
		Hard:   2,
	},
	}
	fmt.Println(printFormat(&p))

	s := St{Dice: &oneroll.DiePool{
		Normal: 2,
		Hard:   4,
	},
	}
	fmt.Println(printFormat(&s))
}

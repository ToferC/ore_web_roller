package main

import (
	"fmt"

	"github.com/toferc/oneroll"
)

func main() {

	c := oneroll.NewReignCharacter("Nornam")

	r1 := oneroll.Roll{
		Actor:  c,
		Action: c.Skills["Athletics"].Name,
	}

	c.Skills["Athletics"].Dice.Normal = 2
	c.Skills["Athletics"].Dice.Expert = 7

	ath := c.Skills["Athletics"]

	athString := ath.FormatDiePool(1)

	fmt.Println(athString)

	r1.Resolve(athString)

	fmt.Println(r1)
}

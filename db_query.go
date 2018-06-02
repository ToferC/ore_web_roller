package main

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
)

// Query and return a Character from DB
func Query(db *pg.DB) {

	c, err := GetCharacter(db)

	if err != nil {
		panic(err)
	}

	for true {
		fmt.Println("\nCharacter Skills:\n")

		for _, v := range c.Skills {
			if oneroll.SkillRated(v) {
				fmt.Println(v)
			}
		}

		skillroll := UserQuery("\nChoose a skill to roll or hit Enter to quit: ")

		if len(skillroll) == 0 {
			fmt.Println("Exiting.")
			break
		}

		validSkill := true

		for k := range c.Skills {
			if skillroll == k {
				validSkill = true
				break
			}
			validSkill = false
		}

		if !validSkill {
			fmt.Println("Not a skill. Try again.")
		} else {

			s := c.Skills[skillroll]

			ds := oneroll.FormSkillDieString(s, 1)

			fmt.Printf("Rolling %s (w/ %s) for %s\n",
				s.Name,
				s.LinkStat.Name,
				c.Name)

			r := oneroll.Roll{
				Actor:  c,
				Action: "act",
			}

			r.Resolve(ds)

			fmt.Println("Rolling!")
			fmt.Println(r)
		}
	}
}

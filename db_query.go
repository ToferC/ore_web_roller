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

QueryActionLoop:
	for true {
		fmt.Println("Would you like to update Statistics or Skills?")

		answer := UserQuery(`
		1: Make a Skill Roll
		2: Mark Damage
		3: Coming Soon...
		4: Coming Soon...

		Or hit Enter to exit: `)

		if len(answer) == 0 {
			fmt.Println("Exiting")
			break QueryActionLoop
		}

		switch answer {
		case "1":
			rollSkill(c)
		case "2":
			//markDamage(db, c)
		default:
			fmt.Println("Not a valid option. Please choose again")
		}
	}
}

func rollSkill(c *oneroll.Character) {

ChooseSkillLoop:
	for true {
		fmt.Println("\nCharacter Skills:\n")

		fmt.Println(oneroll.ShowSkills(c))

		skillroll := UserQuery("\nChoose a skill to roll or hit Enter to quit: ")

		if skillroll == "" {
			fmt.Println("Exiting.")
			break ChooseSkillLoop
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

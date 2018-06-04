package main

import (
	"fmt"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
)

func Update(db *pg.DB) {
	c, err := GetCharacter(db)

	if err != nil {
		panic(err)
	}

UpdateLoop:
	for true {
		fmt.Println("Would you like to update Statistics or Skills?")

		answer := UserQuery(`
    1: Update Statistics
    2: Update Skills
    3: Add a Skill
    4: Delete a Skill

    Or hit Enter to exit: `)

		if len(answer) == 0 {
			fmt.Println("Exiting")
			break UpdateLoop
		}

		switch answer {
		case "1":
			updateStatistics(db, c)
		case "2":
			updateSkills(db, c)
		case "3":
			AddSkill(db, c)
		case "4":
			deleteSkills(db, c)
		default:
			fmt.Println("Not a valid option. Please choose again")
		}
	}
}

func updateStatistics(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Updating Statistics")

	statistics := []*oneroll.Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

UpdateStats:
	for true {

		for i, stat := range statistics {
			fmt.Printf("%d %s\n", i+1, stat)
		}

		fmt.Printf("\nChoose the number for the statistic to update. (1-%d): ", len(statistics))
		fmt.Println("Or hit Enter to exit")

		answer := UserQuery("Your selection: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break UpdateStats
		}

		num, _ := strconv.Atoi(answer)

		if num > 6 || num < 1 {
			fmt.Println("Not a valid statistic. Try again.")
		} else {
			err := updateStat(db, statistics[num-1], c)
			if err != nil {
				panic(err)
			}
			fmt.Println("Updated. Choose another stat or hit Enter to exit.")
		}
	}
}

func updateStat(db *pg.DB, s *oneroll.Statistic, c *oneroll.Character) error {

	fmt.Println(s)

	fmt.Printf("%s has %d normal dice.\n", s.Name, s.Dice.Normal)
	nd := UserQuery("Please enter the new value: ")
	normal, err := strconv.Atoi(nd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Normal = normal
	}

	fmt.Printf("%s has %d hard dice.\n", s.Name, s.Dice.Hard)

	hd := UserQuery("Please enter the new value: ")
	hard, _ := strconv.Atoi(hd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Hard = hard
	}

	fmt.Printf("%s has %d wiggle dice.\n", s.Name, s.Dice.Wiggle)

	wd := UserQuery("Please enter the new value: ")
	wiggle, _ := strconv.Atoi(wd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Wiggle = wiggle
	}

	fmt.Printf("%s has %d spray dice.\n", s.Name, s.Dice.Spray)

	sp := UserQuery("Please enter the new value: ")
	spray, _ := strconv.Atoi(sp)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Spray = spray
	}

	fmt.Printf("%s has %d ranks in go first.\n", s.Name, s.Dice.GoFirst)

	gf := UserQuery("Please enter the new value: ")
	gofirst, _ := strconv.Atoi(gf)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.GoFirst = gofirst
	}

	fmt.Println(c)

	// Update Linked Skills
	for _, skill := range c.Skills {
		if skill.LinkStat.Name == s.Name {
			skill.LinkStat = s
		}
	}

	// Update Willpower
	c.BaseWill = c.Command.Dice.Normal + c.Charm.Dice.Normal
	c.Willpower = c.BaseWill

	// Save character
	err = UpdateCharacter(db, c)
	if err != nil {
		panic(err)
	}

	return err
}

func updateSkills(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Updating Skills")

UpdateSkillsLoop:
	for true {

		fmt.Println(oneroll.ShowSkills(c, true))

		answer := UserQuery("\nType the name of the skill to update or hit Enter to exit: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break UpdateSkillsLoop
		}

		validSkill := false

		for k := range c.Skills {
			if answer == k {
				validSkill = true
				break
			}
			validSkill = false
		}

		if !validSkill {
			fmt.Println("Not a skill. Try again.")

		} else {

			targetSkill := c.Skills[answer]

			err := updateSkill(db, targetSkill, c)
			if err != nil {
				panic(err)
			}
			fmt.Println("Updated. Choose another skill or hit Enter to exit.")
		}
	}
}

func updateSkill(db *pg.DB, s *oneroll.Skill, c *oneroll.Character) error {

	fmt.Println(s)

	if s.ReqSpec {
		fmt.Println("Current specialization is ", s.Specialization)
		spec := UserQuery("Please enter a new specialization or hit Enter to keep the current one: ")
		if len(spec) > 0 {
			s.Specialization = spec
		}
	}

	fmt.Printf("%s has %d normal dice.\n", s.Name, s.Dice.Normal)
	nd := UserQuery("Please enter the new value: ")
	normal, err := strconv.Atoi(nd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Normal = normal
	}

	fmt.Printf("%s has %d hard dice.\n", s.Name, s.Dice.Hard)

	hd := UserQuery("Please enter the new value: ")
	hard, _ := strconv.Atoi(hd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Hard = hard
	}

	fmt.Printf("%s has %d wiggle dice.\n", s.Name, s.Dice.Wiggle)

	wd := UserQuery("Please enter the new value: ")
	wiggle, _ := strconv.Atoi(wd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Wiggle = wiggle
	}

	fmt.Printf("%s has %d spray dice.\n", s.Name, s.Dice.Spray)

	sp := UserQuery("Please enter the new value: ")
	spray, _ := strconv.Atoi(sp)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Spray = spray
	}

	fmt.Printf("%s has %d ranks in go first.\n", s.Name, s.Dice.GoFirst)

	gf := UserQuery("Please enter the new value: ")
	gofirst, _ := strconv.Atoi(gf)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.GoFirst = gofirst
	}

	fmt.Println(c)

	// Save character
	err = UpdateCharacter(db, c)
	if err != nil {
		panic(err)
	}

	return err
}

func AddSkill(db *pg.DB, c *oneroll.Character) {

AddSkillLoop:
	for true {

		fmt.Println(oneroll.ShowSkills(c, true))

		fmt.Println("Adding a new skill")

		s := oneroll.Skill{
			Name: "",
			Dice: &oneroll.DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
			ReqSpec:        false,
			Specialization: "",
		}

		// Get user input for new skill

		answer := UserQuery("Enter the name of the new skill or hit Enter to exit: ")

		if answer == "" {
			break AddSkillLoop
		}

		s.Name = answer

		statistics := []*oneroll.Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

	ChooseStatLoop:
		for true {

			fmt.Println("Statistics:")

			for i, stat := range statistics {
				fmt.Printf("%d - %s\n", i+1, stat)
			}

			fmt.Printf("\nChoose the number for the statistic to update (1-%d) ", len(statistics))
			fmt.Println("or hit Enter to exit")

			answer := UserQuery("\nYour selection: ")

			if answer == "" {
				fmt.Println("Exiting.")
				break ChooseStatLoop
			}

			num, _ := strconv.Atoi(answer)

			if num > 6 || num < 1 {
				fmt.Println("Not a valid statistic. Try again.")
			} else {
				s.LinkStat = statistics[num-1]
				fmt.Println("Updated.")
				break
			}
		}

		sp := UserQuery("Does the skill have a specialization? (Y/N):")

		if sp == "Y" || sp == "y" {
			s.ReqSpec = true
		}

		if s.ReqSpec {
			spec := UserQuery("Enter your specialization: ")
			s.Specialization = spec
		}

		c.Skills[answer] = &s

		updateSkill(db, &s, c)
	}
}

func deleteSkills(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Deleting Skills")

DeleteSkillLoop:
	for true {

		fmt.Println(oneroll.ShowSkills(c, false))

		answer := UserQuery("\nType the name of the skill to delete or hit Enter to exit: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break DeleteSkillLoop
		}

		validSkill := false

		for k := range c.Skills {
			if answer == k {
				validSkill = true
				break
			}
			validSkill = false
		}

		if !validSkill {
			fmt.Println("Not a skill. Try again.")

		} else {

			response := UserQuery("Are you sure you want to delete " + answer + " ? (Y/N)")

			if response == "Y" || response == "y" {
				delete(c.Skills, answer)
				fmt.Println("Deleted.")
			} else {
				fmt.Println("Delete aborted.")
			}

			// Save character
			err := UpdateCharacter(db, c)
			if err != nil {
				panic(err)
			}

		}
		fmt.Println("Deleted.")
	}
}

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
)

func Update(db *pg.DB) {
	c, err := GetCharacter(db)

	if err != nil {
		panic(err)
	}

	for true {
		fmt.Println("Would you like to update Statistics or Skills?")

		answer := UserQuery(`
    1: Update Statistics
    2: Update Skills
    Enter: Exit
    `)

		if len(answer) == 0 {
			fmt.Println("Exiting")
			os.Exit(3)
		}

		switch answer {
		case "1":
			updateStatistics(db, c)
		case "2":
			//updateSkills(db, c)
			fmt.Println("Coming soon.")
		default:
			fmt.Println("Not a valid option. Please choose again")
		}
	}
}

func updateStatistics(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Updating Statistics")

	statistics := []*oneroll.Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

	for true {

		for i, stat := range statistics {
			fmt.Printf("%d - %s", i+1, stat)
		}

		fmt.Printf("\nChoose the number for the statistic to update. (1-%d): ", len(statistics))
		fmt.Println("Or hit Enter to exit")

		answer := UserQuery("Your selection: ")

		if len(answer) == 0 {
			fmt.Println("Exiting.")
			os.Exit(3)
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

	fmt.Printf("%s has %d normal dice.", s.Name, s.Dice.Normal)
	nd := UserQuery("Please enter the new value: ")
	normal, err := strconv.Atoi(nd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Normal = normal
	}

	fmt.Printf("%s has %d hard dice.", s.Name, s.Dice.Hard)

	hd := UserQuery("Please enter the new value: ")
	hard, _ := strconv.Atoi(hd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Hard = hard
	}

	fmt.Printf("%s has %d wiggle dice.", s.Name, s.Dice.Wiggle)

	wd := UserQuery("Please enter the new value: ")
	wiggle, _ := strconv.Atoi(wd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Wiggle = wiggle
	}

	fmt.Printf("%s has %d spray dice.", s.Name, s.Dice.Spray)

	sp := UserQuery("Please enter the new value: ")
	spray, _ := strconv.Atoi(sp)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Spray = spray
	}

	fmt.Printf("%s has %d ranks in go first.", s.Name, s.Dice.GoFirst)

	gf := UserQuery("Please enter the new value: ")
	gofirst, _ := strconv.Atoi(gf)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.GoFirst = gofirst
	}

	c.Display()

	for _, skill := range c.Skills {
		if skill.LinkStat.Name == s.Name {
			skill.LinkStat = s
		}
	}

	// Save character
	err = UpdateCharacter(db, c)
	if err != nil {
		panic(err)
	}

	return err
}

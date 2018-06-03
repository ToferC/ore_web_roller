package oneroll

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Max returns the larger of two ints
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// RollDie rolls and sum dice
func RollDie(max, min, numDice int) int {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	result := 0
	for i := 1; i < numDice+1; i++ {
		roll := r1.Intn(max+1-min) + min
		result += roll
	}
	return result
}

// TrimSliceBrackets trims the brackets from a slice and return ints as a string
func TrimSliceBrackets(s []int) string {
	rs := fmt.Sprintf("%d", s)
	rs = strings.Trim(rs, "[]")
	return rs
}

// ParseNumRolls checks how many die rolls are required
func ParseNumRolls(s string) (int, error) {

	re := regexp.MustCompile("[0-9]+")

	var num int
	var numString string

	numString = re.FindString(s)
	num, err := strconv.Atoi(numString)
	if err != nil {
		num = 1
	}
	return num, err
}

// SkillRated returns true if a skill has any points in it
func SkillRated(s *Skill) bool {
	if s.Dice.Normal+s.Dice.Hard+s.Dice.Wiggle+s.Dice.Spray+s.Dice.GoFirst > 0 {
		return true
	}
	return false
}

// ShowSkills shows skills grouped under stats
// all bool determines if all skills are shown or just the ones with dice in them.
func ShowSkills(c *Character, all bool) string {
	statistics := []*Statistic{c.Body, c.Coordination, c.Sense, c.Mind, c.Command, c.Charm}

	var text string

	for _, stat := range statistics {
		text += fmt.Sprintf("%s\n", stat)
		for _, skill := range c.Skills {
			if skill.LinkStat.Name == stat.Name {
				if all {
					text += fmt.Sprintf("-- %s\n", skill)
				} else {
					if SkillRated(skill) {
						text += fmt.Sprintf("-- %s\n", skill)
					}
				}
			}
		}
	}
	return text
}

// UserQuery creates and question and returns the User's input as a string
func UserQuery(q string) string {
	question := bufio.NewReader(os.Stdin)
	fmt.Print(q)
	r, _ := question.ReadString('\n')

	input := strings.Trim(r, " \n")

	return input
}

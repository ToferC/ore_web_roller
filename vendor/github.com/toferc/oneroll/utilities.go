package oneroll

import (
	"fmt"
	"math/rand"
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

// VerifyLessThan10 checks and reduces die pools to less than 10d
func VerifyLessThan10(nd, hd, wd int) (int, int, int) {

	if nd+hd+wd > 10 {

		fmt.Println("Error: Can't roll more than 10 dice. Reducing to less than 10.")
		fmt.Printf(fmt.Sprintf("Current Dice: %dd+%dhd+%dwd.\n", nd, hd, wd))

		// Remove normal dice first
		for nd > 0 && nd+hd+wd > 10 {
			fmt.Printf("reduced Normal dice from %d to %d. \n", nd, nd-1)
			nd--
			fmt.Printf(fmt.Sprintf("Current Dice: %dd+%dhd+%dwd.\n", nd, hd, wd))
		}

		// Reduce hard dice next
		for hd > 0 && nd+hd+wd > 10 {
			fmt.Printf("reduced Hard dice from %d to %d. \n", hd, hd-1)
			hd--
			fmt.Printf(fmt.Sprintf("Current Dice: %dd+%dhd+%dwd.\n", nd, hd, wd))
		}

		// Reduce wiggle dice last
		for wd > 0 && nd+hd+wd > 10 {
			fmt.Printf("reduced Wiggle dice from %d to %d. \n", wd, wd-1)
			wd--
			fmt.Printf(fmt.Sprintf("Current Dice: %dd+%dhd+%dwd.\n", nd, hd, wd))
		}

		return nd, hd, wd

	}

	return nd, hd, wd
}

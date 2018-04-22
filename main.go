package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/toferc/oneroll"
)

func main() {

	c := oneroll.NewCharacter("Baron")
	c.Display()

	d := oneroll.NewCharacter("Duke")
	d.Display()

	fmt.Println("Let the Arm Wrestingling Commence!")

	actingSkill := c.Skills["Athletics"]

	actingSkill.Dice.Spray = 1

	s := oneroll.FormSkillDieString(actingSkill, 2)

	fmt.Printf("Rolling Athletics (Body+Athletics) for %s", c.Name)
	r := oneroll.Roll{
		Actor:  c,
		Action: "act",
	}

	r.Resolve(s)

	fmt.Println("Rolling!")
	fmt.Println(r)

	opposingSkill := d.Skills["Athletics"]
	s2 := oneroll.FormSkillDieString(opposingSkill, 1)

	fmt.Printf("Rolling Athletics (Body+Athletics) for %s\n", d.Name)
	r2 := oneroll.Roll{
		Actor:  d,
		Action: "oppose",
	}

	r2.Resolve(s2)

	fmt.Println("Rolling!")
	fmt.Println(r2)

	oneroll.OpposedRoll(&r, &r2)

	fmt.Println("Starting Webserver at port 8080")
	http.HandleFunc("/roll/", RollHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

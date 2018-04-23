package oneroll

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/andlabs/ui"
)

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

// GUI renders a GUI for the die rolling app
func GUI() {
	err := ui.Main(func() {
		ndInput := ui.NewEntry()
		hdInput := ui.NewEntry()
		wdInput := ui.NewEntry()
		numInput := ui.NewEntry()
		goFirstInput := ui.NewEntry()
		sprayInput := ui.NewEntry()
		acInput := ui.NewEntry()
		button := ui.NewButton("Roll")
		results := ui.NewLabel("")
		box := ui.NewVerticalBox()
		rowD := ui.NewHorizontalBox()
		rowHD := ui.NewHorizontalBox()
		rowWD := ui.NewHorizontalBox()
		rowGF := ui.NewHorizontalBox()
		rowSp := ui.NewHorizontalBox()
		rowNum := ui.NewHorizontalBox()
		rowAc := ui.NewHorizontalBox()

		box.Append(ui.NewLabel("Build your ORE Dice Pool\n"), false)
		box.Append(ui.NewHorizontalSeparator(), false)

		rowD.Append(ndInput, false)
		rowD.Append(ui.NewLabel(" Normal Dice"), false)
		rowHD.Append(hdInput, false)
		rowHD.Append(ui.NewLabel(" Hard Dice"), false)
		rowWD.Append(wdInput, false)
		rowWD.Append(ui.NewLabel(" Wiggle Dice"), false)

		box.Append(rowD, false)
		box.Append(rowHD, false)
		box.Append(rowWD, false)
		box.Append(ui.NewHorizontalSeparator(), false)
		box.Append(ui.NewLabel(""), false)

		rowGF.Append(goFirstInput, false)
		rowGF.Append(ui.NewLabel(" Go First Levels"), false)

		rowSp.Append(sprayInput, false)
		rowSp.Append(ui.NewLabel(" Spray Levels"), false)

		rowNum.Append(numInput, false)
		rowNum.Append(ui.NewLabel(" Number of Rolls"), false)

		rowAc.Append(acInput, false)
		rowAc.Append(ui.NewLabel(" Number of Actions"), false)

		box.Append(ui.NewHorizontalSeparator(), false)
		box.Append(rowGF, false)
		box.Append(rowSp, false)
		box.Append(rowAc, false)
		box.Append(rowNum, false)
		box.Append(button, false)
		box.Append(results, false)

		window := ui.NewWindow("ORE Die Roller", 300, 600, false)
		window.SetMargined(true)
		window.SetChild(box)

		button.OnClicked(func(*ui.Button) {

			var resultString string

			c := Character{
				Name: "Player",
			}

			numRolls, err := ParseNumRolls(numInput.Text())

			if err != nil {
				resultString += "Invalid number of rolls. Set to 1.\n\n"
			}

			for x := 1; x < numRolls+1; x++ {

				resultString += fmt.Sprintf("Roll #%d\n", x)

				roll := Roll{
					Actor:  &c,
					Action: "Act",
				}

				text := fmt.Sprintf("%sac+%sd+%shd+%swd+%sgf+%ssp",
					acInput.Text(),
					ndInput.Text(),
					hdInput.Text(),
					wdInput.Text(),
					goFirstInput.Text(),
					sprayInput.Text(),
				)

				r, err := roll.Resolve(text)

				if err != nil {
					resultString += fmt.Sprintf("%s", err)
				} else {
					resultString += fmt.Sprintf("%s", r)
				}
			}
			results.SetText(fmt.Sprintf("%s", resultString))

		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

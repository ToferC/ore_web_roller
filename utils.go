package main

import (
	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
	"github.com/toferc/ore_web_roller/models"
)

type Element interface {
	ChooseDiePool() error
}

func ChooseDice(e Element) error {
	err := e.ChooseDiePool()
	if err != nil {
		return err
	}
	return nil
}

func convertToModels(*pg.DB) error {

	//Add Checks before running

	var chars []*oneroll.Character
	var pows map[string]oneroll.Power

	author, err := database.PKLoadUser(db, int64(1))
	if err == nil {
		chars, _ = database.ListCharacters(db)
		pows, _ = database.ListPowers(db)

		for _, c := range chars {
			tCM := models.CharacterModel{
				Author:    author,
				Character: c,
				Open:      true,
			}
			database.SaveCharacterModel(db, &tCM)
			database.DeleteCharacter(db, c.ID)
		}

		for _, v := range pows {
			tPM := models.PowerModel{
				Author: author,
				Power:  &v,
				Open:   true,
			}
			database.SavePowerModel(db, &tPM)
			database.DeletePower(db, v.ID)
		}
		return nil
	}
	return err
}

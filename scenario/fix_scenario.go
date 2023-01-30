package scenario

import (
	"github.com/PDC-Repository/newauth/newauth/models"
	"gorm.io/gorm"
)

func NewBot(db *gorm.DB) (*models.Bot, func()) {
	bot := models.Bot{
		Name: "bot new test",
		Desc: "Description",
	}
	err := db.Save(&bot).Error
	if err != nil {
		panic(err)
	}
	return &bot, func() {
		err := db.Delete(&bot).Error
		if err != nil {
			panic(err)
		}
	}
}

func NewTeam(db *gorm.DB) (*models.Team, func()) {
	team := models.Team{
		Name:        "team test api",
		Description: "asdasdasd",
	}
	err := db.Save(&team).Error
	if err != nil {
		panic(err)
	}

	return &team, func() {
		db.Delete(&team)
	}
}

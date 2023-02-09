package scenario

import (
	"log"
	"time"

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

func NewUser(db *gorm.DB) (*models.User, func()) {
	idnya := generateUsername()

	tgl := time.Now().AddDate(0, -1, 0)

	user := models.User{
		Name:      idnya,
		Email:     idnya + "@gmail.com",
		Username:  idnya,
		LastReset: tgl,
	}
	user.SetPassword("password")

	err := db.Create(&user).Error
	if err != nil {
		log.Panicln("gagal create user")
	}

	return &user, func() {
		err := db.Unscoped().Delete(&user, user.ID).Error
		if err != nil {
			log.Println("gagal delete")
		}
	}
}

func NewBotToken(db *gorm.DB) (*models.User, *models.BotToken, *models.Bot, string, func()) {
	user, Tuser := NewUser(db)
	bot, Tbot := NewBot(db)
	team, Tteam := NewTeam(db)
	password := "password"

	botToken := models.BotToken{
		BotID:  bot.ID,
		UserID: user.ID,
		TeamID: team.ID,
	}
	botToken.SetPwd(password)

	err := db.Create(&botToken).Error
	if err != nil {
		panic(err)
	}

	return user, &botToken, bot, password, func() {
		db.Delete(&botToken)
		Tuser()
		Tbot()
		Tteam()
	}
}

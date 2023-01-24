package newauth

import (
	"log"
	"sync"

	"github.com/PDC-Repository/newauth/config"
	"github.com/PDC-Repository/newauth/newauth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbOnce sync.Once

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.UserTeam{},
	)

	if err != nil {
		panic(err)
	}
}

func NewDatabase() *gorm.DB {
	dbOnce.Do(func() {
		log.Println("initialize database")
		db, err := gorm.Open(postgres.Open(config.DatabaseUri), &gorm.Config{})

		if err != nil {
			panic(err)
		}
		AutoMigrate(db)
		DB = db
	})

	return DB
}

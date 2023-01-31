package newauth

import (
	"log"
	"sync"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
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
		&models.Bot{},
		&models.Quota{},
		&models.BotToken{},
	)

	if err != nil {
		panic(err)
	}
}

func NewDatabase() *gorm.DB {
	dbOnce.Do(func() {
		log.Println("initialize database")
		var db *gorm.DB
		var err error

		if config.Config.DevMode {
			db, err = gorm.Open(postgres.Open(config.Config.Database.CreateDsn()), &gorm.Config{})
		} else {
			db, err = gorm.Open(postgres.New(
				postgres.Config{
					DriverName: "cloudsqlpostgres",
					DSN:        config.Config.Database.CreateDsn(),
				}),
			)
		}

		if err != nil {
			panic(err)
		}
		AutoMigrate(db)
		DB = db
	})

	return DB
}

package newauth

import (
	"log"
	"sync"

	"github.com/PDC-Repository/newauth/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbOnce sync.Once

func NewDatabase() *gorm.DB {

	dbOnce.Do(func() {
		log.Println("initialize database")
		db, err := gorm.Open(postgres.Open(config.DatabaseUri), &gorm.Config{})

		if err != nil {
			panic(err)
		}

		DB = db
	})

	return DB
}

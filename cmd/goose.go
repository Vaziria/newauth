package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/PDC-Repository/newauth/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	flags    = flag.NewFlagSet("goose", flag.ExitOnError)
	dbstring = flags.String("dbstring", "user=admin password=admin dbname=postgres sslmode=disable", "connection string")
	dir      = flags.String("dir", "./migrations", "directory with migration files")
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}
	command := args[0]

	db, err := goose.OpenDBWithDriver("postgres", *dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}

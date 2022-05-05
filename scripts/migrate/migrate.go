package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/Shelex/split-specs-v2/env"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config := env.ReadEnv()

	var direction string
	var version uint
	flag.StringVar(&direction, "direction", "up", "Apply up or down migrations")
	flag.UintVar(&version, "version", 0, "Apply N 'up' migrations")
	flag.Parse()

	log.Printf("migrating: direction:%s; version:%d", direction, version)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sourceURL := "file://" + dir + "/repository/migrations"
	dbURL := config.DbConnectionUrl + "&sslmode=disable"

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	if version != 0 {
		err = m.Migrate(version)
	} else {
		if direction == "up" {
			err = m.Up()
		} else {
			err = m.Down()
		}
	}
	if err != nil {
		panic(err)
	}

	version, _, err = m.Version()
	log.Println("current version is " + strconv.FormatUint(uint64(version), 10))
	if err != nil {
		panic(err)
	}
}

package main

import (
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

func initCommand(c *cli.Context) error {
	conf.ContentPgAddr = c.GlobalString("pg-addr")
	conf.ContentPgUser = c.GlobalString("pg-user")
	conf.ContentPgPass = c.GlobalString("pg-pass")
	conf.ContentPgDB = c.GlobalString("pg-db")

	err := connectDB()
	if err != nil {
		log.Warn(err)
		return nil
	}

	err = executeInit()
	if err != nil {
		log.Warn(err)
		return nil
	}
	log.Info("successfully created db tables")
	return nil
}

func executeInit() error {
	for _, stmt := range getMigrations(false) {
		_, err := contentdb.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func getMigrations(invert bool) []string {
	migrations, err := AssetDir("migrations/pg")
	if err != nil {
		log.Fatal(err)
	}

	// order.
	sort.Strings(migrations)

	if invert {
		migrations = reverse(migrations)
	}

	result := []string{}
	for _, migration := range migrations {
		if strings.Contains(migration, "down") {
			continue
		}
		sql, err := Asset("migrations/pg/" + migration)
		if err != nil {
			log.Fatal(err)
		}

		stmt := string(sql)
		if stmt == "" {
			continue
		}

		result = append(result, stmt)
	}
	return result
}

func reverse(files []string) []string {
	for i := 0; i < len(files)/2; i++ {
		j := len(files) - i - 1
		files[i], files[j] = files[j], files[i]
	}
	return files
}

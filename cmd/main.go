package main

import (
	"os"
	"sort"

	"github.com/piquette/edgr/database"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var (
	// content database.
	contentdb *database.Handle
	connected bool
)

func main() {
	app := cli.NewApp()
	app.Name = "edgr"
	app.Version = "0.0.1"
	app.Usage = "Retrieve and store SEC filings for corporations"
	app.UsageText = "edgr [global flags] COMMAND [command flags]"

	app.Flags = buildGlobalFlags()
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Initializes a postgres database that can store SEC data",
			Action: initCommand,
		},
		{
			Name:   "get",
			Usage:  "Retrieves and stores SEC filings",
			Flags:  buildGetFlags(),
			Action: getCommand,
		},
		{
			Name:  "filers",
			Usage: "Manage the universe of entities that file with the SEC",
			Subcommands: []cli.Command{
				{
					Name:   "init",
					Usage:  "Retrieves and stores any filers that can reasonably be matched to a publicly traded stock symbol",
					Action: filersInitCommand,
					Flags:  buildFilersFlags(),
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	defer disconnectDB()
	_ = app.Run(os.Args)
}

func connectDB(c *cli.Context) error {
	// Connect to Postgres db.
	log.Info("connecting to postgres..")
	copts := database.Options{
		Addr:     c.GlobalString("pg-addr"),
		User:     c.GlobalString("pg-user"),
		Password: c.GlobalString("pg-pass"),
		Database: c.GlobalString("pg-db"),
	}

	contentdb = database.Open(copts)
	_, err := contentdb.Exec("SELECT NULL")
	if err != nil {
		return err
	}
	log.Info("connected")
	connected = true
	return nil
}

func disconnectDB() {
	if !connected {
		return
	}
	contentdb.Close()
	log.Info("disconnected from postgres")
}

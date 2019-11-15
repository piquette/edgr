package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/piquette/edgr/database"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

var connected bool

// Edgr holds the external services needed to accomplish tasks.
type Edgr struct {
	FilingDao database.FilingDao
	FilerDao  database.FilerDao
}

// Config is the configuration.
var conf = struct {
	// PostgreSQL/persistence settings
	ContentPgAddr string
	ContentPgUser string
	ContentPgPass string
	ContentPgDB   string
	// Settings.
	LetterStart string
	SymbolStart string
	FillMode    string
	StopDate    string
}{}

var (
	dirRegex = regexp.MustCompile(`<td><a href="(.*?)"><img`)
	urlRegex = regexp.MustCompile(`.*<a href="(.*?)index.html"><img`)
	// content database.
	contentdb *database.Handle
)

func main() {
	app := cli.NewApp()
	app.Name = "edgr"
	app.Version = "0.0.1"
	app.Usage = "Retrieve and store SEC filings for corporations"
	app.UsageText = "edgr [global flags] COMMAND [command flags]"

	app.Flags = buildFlags()
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Initializes a postgres database that can store SEC data",
			Action: initCommand,
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	defer disconnectDB()
	_ = app.Run(os.Args)
}

func connectDB() error {

	// Connect to Postgres db.
	log.Info("connecting to db")
	copts := database.Options{
		Addr:     conf.ContentPgAddr,
		User:     conf.ContentPgUser,
		Password: conf.ContentPgPass,
		Database: conf.ContentPgDB,
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
	log.Info("disconnected from db")

}

func buildFlags() []cli.Flag {
	// 	LetterStart string `flag:"letter" env:"LETTER" usage:"LETTER"`
	// 	SymbolStart string `flag:"start" env:"START" usage:"START"`
	// 	FillMode    string `flag:"fill-mode" env:"FILL_MODE" usage:"Fill mode"`
	// 	StopDate    string `flag:"stop-date" env:"STOP_DATE" usage:"Stop date YYYY-MM-DD"`
	// 	LetterStart:   "A",
	// 	SymbolStart:   "",
	// 	FillMode:      "alpha",
	// 	StopDate:      "2019-06-01",
	return []cli.Flag{
		cli.StringFlag{
			Name:     "pg-addr",
			Value:    "localhost:5432",
			Usage:    "PostgreSQL ~address~",
			EnvVar:   "PG_ADDR",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "pg-user",
			Value:    "postgres",
			Usage:    "PostgreSQL ~username~",
			EnvVar:   "PG_USER",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "pg-pass",
			Value:    "postgres",
			Usage:    "PostgreSQL ~password~",
			EnvVar:   "PG_PASS",
			FilePath: "EDGRFILE",
		},
		cli.StringFlag{
			Name:     "pg-db",
			Value:    "postgres",
			Usage:    "PostgreSQL ~database~",
			EnvVar:   "PG_DB",
			FilePath: "EDGRFILE",
		},
	}
}

// findListURLs parses the list of idx urls out of the directory page.
func findListURLs(html string) []string {
	matches := dirRegex.FindAllStringSubmatch(html, -1)
	if matches == nil || len(matches) == 1 {
		fmt.Println("NO MATCHES")
		return nil
	}

	urls := []string{}
	for _, m := range matches {
		urls = append(urls, "https://sec.gov"+m[1])
	}
	return urls
}

// findIdxURL parses the text document url out of the index page.
func findIdxURL(html string) string {
	matches := urlRegex.FindStringSubmatch(html)
	if matches == nil || len(matches) == 1 {
		fmt.Println("NO MATCHES")
		return ""
	}
	return "https://sec.gov" + matches[1] + "index.html"
}

package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/piquette/edgr/database"
)

// Edgr holds the external services needed to accomplish tasks.
type Edgr struct {
	FilingDao database.FilingDao
	FilerDao  database.FilerDao
}

var conf = struct {
	// PostgreSQL/persistence settings
	ContentPgAddr string `flag:"content-pg-addr" usage:"PostgreSQL ~address~" env:"CONTENT_PG_ADDR"`
	ContentPgUser string `flag:"content-pg-user" usage:"PostgreSQL ~username~" env:"CONTENT_PG_USER"`
	ContentPgPass string `flag:"content-pg-pass" usage:"PostgreSQL ~password~" env:"CONTENT_PG_PASS"`
	ContentPgDB   string `flag:"content-pg-db" usage:"PostgreSQL ~database~" env:"CONTENT_PG_DB"`
	// Settings.
	LetterStart string `flag:"letter" env:"LETTER" usage:"LETTER"`
	SymbolStart string `flag:"start" env:"START" usage:"START"`
	FillMode    string `flag:"fill-mode" env:"FILL_MODE" usage:"Fill mode"`
	StopDate    string `flag:"stop-date" env:"STOP_DATE" usage:"Stop date YYYY-MM-DD"`
}{
	ContentPgAddr: "localhost:5432",
	ContentPgUser: "postgres",
	ContentPgPass: "postgres",
	ContentPgDB:   "postgres",
	LetterStart:   "A",
	SymbolStart:   "",
	FillMode:      "alpha",
	StopDate:      "2019-06-01",
}

var (
	dirRegex = regexp.MustCompile(`<td><a href="(.*?)"><img`)
	urlRegex = regexp.MustCompile(`.*<a href="(.*?)index.html"><img`)
	// content database.
	contentdb *database.Handle
)

func main() {

	//letter := strings.Split(, ",")
	log.Println("starting backfill for letter:", conf.LetterStart, "symbol:", conf.SymbolStart)

	// Connect to Postgres db.
	log.Println("connecting to db")
	// Connect to Postgres db.
	log.Println("connecting to content-database")
	copts := database.Options{
		Addr:     conf.ContentPgAddr,
		User:     conf.ContentPgUser,
		Password: conf.ContentPgPass,
		Database: conf.ContentPgDB,
	}

	contentdb = database.Open(copts)
	defer contentdb.Close()

	_, err := contentdb.Exec("SELECT NULL")
	if err != nil {
		log.Fatalln(err)
	}

	e := &Edgr{
		FilingDao: contentdb.NewFilingDao(),
	}

	if conf.FillMode == "alpha" {
		e.executeAlpha()
		return
	}

	if conf.FillMode == "time" {
		e.executeTime()
		return
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

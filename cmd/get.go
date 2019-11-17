package main

import (
	"fmt"
	"time"

	"github.com/piquette/edgr/core/model"
	"github.com/piquette/edgr/database"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func getCommand(c *cli.Context) error {
	return nil
}

func executeAlpha(dao database.FilerDao) {
	//
	found, filers, err := dao.GetSet(conf.LetterStart)
	if err != nil {
		log.Fatal(err)
	}
	if !found {
		return
	}

	hitSymbol := false
	fdao := contentdb.NewFilingDao()
	for _, filer := range filers {
		if conf.SymbolStart != "" && filer.Symbol != conf.SymbolStart && !hitSymbol {
			continue
		}
		hitSymbol = true
		backfillFilerWithCutoff(fdao, *filer, nil)
	}
}

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func executeTime(dao database.FilerDao) {

	cutoffTime, err := time.Parse("2006-01-02", conf.StopDate)
	if err != nil {
		log.Fatal(err)
	}

	for _, letter := range alphabet {
		found, filers, err := dao.GetSet(letter)
		if err != nil {
			log.Fatal(err)
		}
		if !found {
			return
		}

		dao := contentdb.NewFilingDao()
		for _, filer := range filers {
			backfillFilerWithCutoff(dao, *filer, &cutoffTime)
		}
	}
}

// backfillFilerWithCutoff backfills a filer with a time cutoff.
func backfillFilerWithCutoff(dao database.FilingDao, filer model.Filer, cutoffTime *time.Time) {
	dirPage, err := getPage("https://www.sec.gov/Archives/edgar/data/"+filer.CIK, 2)
	if err != nil {
		log.Println(err)
		return
	}

	urls := findListURLs(dirPage)

	for _, u := range urls {
		docsPage, err := getPage(u, 2)
		if err != nil {
			log.Println(err)
			continue
		}

		idxURL := findIdxURL(docsPage)
		if idxURL == "" {
			log.Println("couldnt find idx url")
			continue
		}

		filing, err := buildFiling(filer, idxURL)
		if err != nil {
			log.Println(err)
			continue
		}

		if cutoffTime != nil {
			// check cutoff time.
			if filing.Filing.EdgarTime.Before(*cutoffTime) {
				log.Println("hit cutoff time")
				return
			}
		}

		// Do stuff with the filing...
		fmt.Println(filing.Filing.Symbol, " - ", filing.Filing.Accession)

		// Match symbols.
		symbols := []string{}
		// for _, cik := range filing.Filing.AllCIKs {
		// 	sym, err := b.Redis.GetSymbol(cik)
		// 	if err != nil || sym == "" {
		// 		continue
		// 	}
		// 	symbols = append(symbols, sym)
		// }
		filing.Filing.AllSymbols = symbols

		created, exists, err := dao.Add(filing.Filing)
		if err != nil {
			log.Println(err)
			continue
		}
		if exists {
			// get id.
			_, fid, _ := dao.Get(filing.Filing.Accession)
			filing.Filing.ID = fid.ID

			_, err = dao.Update(filing.Filing)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(filing.Filing.FilerRelation)
			log.Println("already existed, updated")
			continue
		}

		err = dao.AddDocuments(created, filing.Docs)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("Success")
	}
}

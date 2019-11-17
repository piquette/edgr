package main

import (
	"fmt"
	"math"

	"github.com/schollz/progressbar/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/piquette/edgr/core"
	"github.com/piquette/edgr/database"
)

var (
	iexSymbolsURL = "https://api.iextrading.com/1.0/ref-data/symbols?format=csv"
	secCompanyURL = "https://www.sec.gov/cgi-bin/browse-edgar?action=getcompany&CIK=%s&start=0&count=1&output=atom"
	iexCompanyURL = "https://api.iextrading.com/1.0/stock/spy/company"
)

func filersInitCommand(c *cli.Context) error {
	// Connect to db.
	err := connectDB(c)
	if err != nil {
		log.Warn(err)
		return nil
	}

	// TODO: clean up this mess..
	list := []core.Company{}
	getAll := c.Bool("all")
	sicGroup := c.String("sic")

	if !getAll && sicGroup == "" {
		// try single symbol.
		singleSymbol := c.String("symbol")
		if singleSymbol == "" {
			log.Info("please specify a filer to fetch")
			return nil
		}
		log.Info("retrieving single company")
		list = append(list, core.Company{Symbol: singleSymbol})

	} else {
		// Get all.
		log.Info("retrieving list of companies..")
		log.Info("gathering filer information, this will take a long time..")
		list, err := core.GetPublicCompanies()
		if err != nil {
			return err
		}
		if len(list) == 0 {
			return fmt.Errorf("could not find companies")
		}
	}

	err = storeFilers(list, contentdb.NewFilerDao(), sicGroup)
	if err != nil {
		log.Warn(err)
		return nil
	}
	log.Info("finished")
	return nil
}

func storeFilers(list []core.Company, dao database.FilerDao, sicGroup string) error {
	// match symbols.
	bar := progressbar.NewOptions(100, progressbar.OptionSetPredictTime(true))
	totalSize := float64(len(list))
	prog := 0
	for i, c := range list {

		// Find info for that symbol.
		filer, cikErr := core.GetFiler(c.Symbol)
		if cikErr != nil {
			// Handle.
			continue
		}

		if sicGroup != "" {
			if filer.SIC != sicGroup {
				continue
			}
		}

		_, dbErr := dao.Put(filer)
		if dbErr != nil {
			// Handle.
			log.Warn(dbErr)
			continue
		}
		val := int(math.Ceil((float64(i) / totalSize) * 100))
		if val > prog {
			bar.Add(val)
			prog = val
		}
	}
	bar.Finish()
	return nil
}

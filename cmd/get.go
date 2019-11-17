package main

import (
	"github.com/piquette/edgr/core"
	"github.com/piquette/edgr/core/model"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func getCommand(c *cli.Context) error {
	// symbol
	// sic
	// form
	// stop
	symbol := c.String("symbol")
	sic := c.String("sic")

	if sic == "" && symbol == "" {
		log.Info("please specify a filer")
		return nil
	}

	// Connect to db.
	err := connectDB(c)
	if err != nil {
		log.Warn(err)
		return nil
	}

	// Get filers based on symbol or sic.
	var filers []*model.Filer
	var found bool
	filerDAO := contentdb.NewFilerDao()

	if symbol != "" {
		// symbol.
		found, filers, err = filerDAO.GetSet(symbol)
	} else {
		// sic.
		found, filers, err = filerDAO.GetSetBySIC(sic)
	}

	if err != nil {
		log.Warn(err)
		return nil
	}
	if !found {
		log.Info("could not find matching filer(s) in db")
		return nil
	}

	// Cycle through filers and request filings based on form+stoptime.
	filingDAO := contentdb.NewFilingDao()
	for _, filer := range filers {

		filings, err := core.GetFilings(filer.CIK, c.String("form"), c.String("stop"))
		if err != nil {
			log.Warn(err)
			return nil
		}

		// Store filings.
		for _, secfiling := range filings {
			secfiling.Filing.Filer = filer.Name
			secfiling.Filing.Symbol = filer.Symbol

			created, exists, err := filingDAO.Add(secfiling.Filing)
			if err != nil {
				log.Warn(err)
				continue
			}
			if exists {
				// get id.
				_, fid, _ := filingDAO.Get(secfiling.Filing.Accession)
				secfiling.Filing.ID = fid.ID

				_, err = filingDAO.Update(secfiling.Filing)
				if err != nil {
					log.Warn(err)
					continue
				}
				log.Info(secfiling.Filing.FilerRelation)
				log.Info("already existed, updated")
				continue
			}
			err = filingDAO.AddDocuments(created, secfiling.Docs)
			if err != nil {
				log.Warn(err)
				continue
			}
		}
	}
	return nil
}

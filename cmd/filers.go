package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/piquette/edgr/core/model"
	"golang.org/x/net/html/charset"
)

var (
	iexSymbolsURL = "https://api.iextrading.com/1.0/ref-data/symbols?format=csv"
	secCompanyURL = "https://www.sec.gov/cgi-bin/browse-edgar?action=getcompany&CIK=%s&start=0&count=1&output=atom"
	iexCompanyURL = "https://api.iextrading.com/1.0/stock/spy/company"
)

// RssFeed is the feed obj.
type RssFeed struct {
	Info CompanyInfo `xml:"company-info"`
}

// CompanyInfo internal to rssfeed obj.
type CompanyInfo struct {
	CIK     string `xml:"cik"`
	SIC     string `xml:"assigned-sic,omitempty"`
	SICDesc string `xml:"assigned-sic-desc,omitempty"`
	Name    string `xml:"conformed-name"`
}

// GetFilers runs the period fetch job.
func (s *Edgr) GetFilers() {
	log.Println("executing job")
	// Get symbols.
	// TODO: tag any newly matched ciks with their symbols in filings table

	companyTable, err := fetchCSV(iexSymbolsURL)
	if err != nil {
		// handle with email.
		log.Println(err)
		return
	}
	if len(companyTable) == 0 {
		// handle with email.
		return
	}

	// match symbols.
	newCompanies := map[string]string{}
	log.Println("got table, now matching symbols")

	for _, tRow := range companyTable {

		// Symbol.
		symbol := tRow[0]

		// Find info+cik for that symbol.
		log.Println(symbol)
		info, cikErr := getCompanyInfo(symbol)
		if cikErr != nil {
			// Handle.
			//log.Println(cikErr)
			continue
		}

		// Post to cache.
		// exists, cacheErr := s.SymbolsCache.SetPair(info.CIK, symbol)
		// if cacheErr != nil {
		// 	// Handle.
		// 	log.Println(cacheErr)
		// 	continue
		// }
		// if exists {
		// 	// Handle.
		// 	//log.Printf("%s = %s already exists\n", info.CIK, symbol)
		// 	continue
		// }

		// try to fetch company info..
		// maybe later.. iex has a company url.

		// Add to db.
		filer := &model.Filer{
			CIK:            info.CIK,
			Symbol:         symbol,
			SIC:            info.SIC,
			SICDescription: info.SICDesc,
			Name:           info.Name,
		}
		_, dbErr := s.FilerDao.Put(filer)
		if dbErr != nil {
			// Handle.
			log.Println(dbErr)
			continue
		}

		// Add to new companies.
		newCompanies[symbol] = info.Name
	}
	printSuccess(newCompanies)
}

func bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func printSuccess(companies map[string]string) {
	log.Println("Added the following companies:")
	log.Println()

	for sym, nme := range companies {
		log.Println(" -", bold(sym+":"), " (-"+nme+")")
	}

	log.Println()
	log.Println("Done.")
	log.Println()
}

func getCompanyInfo(symbol string) (info CompanyInfo, err error) {
	// get the cik for each symbol.
	// tedious process...
	url := fmt.Sprintf(secCompanyURL, symbol)

	httpClient := http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var feed RssFeed
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&feed)
	if err != nil {
		return
	}

	if feed.Info.CIK == "" {
		err = fmt.Errorf("no cik found in response data")
		return
	}
	if feed.Info.Name == "" {
		err = fmt.Errorf("no name found in response data")
		return
	}

	return feed.Info, nil
}

func fetchCSV(url string) (table [][]string, err error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)
	r.FieldsPerRecord = -1
	table, err = r.ReadAll()
	return
}

package core

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
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

// Company is a simple struct for a single company.
type Company struct {
	Name   string
	Symbol string
}

// rssFeed is the feed obj.
type rssFeed struct {
	Info secFilerInfo `xml:"company-info"`
}

type secFilerInfo struct {
	CIK     string `xml:"cik"`
	SIC     string `xml:"assigned-sic,omitempty"`
	SICDesc string `xml:"assigned-sic-desc,omitempty"`
	Name    string `xml:"conformed-name"`
}

// GetPublicCompanies returns a list of public companies.
func GetPublicCompanies() ([]Company, error) {

	req, err := http.NewRequest("GET", iexSymbolsURL, nil)
	if err != nil {
		return []Company{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Company{}, err
	}
	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)
	r.FieldsPerRecord = -1
	table, err := r.ReadAll()
	result := []Company{}
	for _, row := range table {
		sym := row[0]
		nme := row[1]
		if len(sym) > 5 || nme == "" {
			continue
		}
		result = append(result, Company{
			Symbol: sym,
			Name:   nme,
		})
	}
	return result, nil
}

// GetFiler gets a single filer from the SEC website based on symbol.
func GetFiler(symbol string) (filer *model.Filer, err error) {
	// get the cik for each symbol.
	// tedious process...
	url := fmt.Sprintf(secCompanyURL, symbol)

	httpClient := http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var feed rssFeed
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

	return &model.Filer{
		CIK:            feed.Info.CIK,
		Symbol:         symbol,
		SIC:            feed.Info.SIC,
		SICDescription: feed.Info.SICDesc,
		Name:           feed.Info.Name,
	}, nil
}

package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/piquette/edgr/core/documents"

	"github.com/piquette/edgr/core/model"
)

// ArchivedFiling is cool.
type ArchivedFiling struct {
	Filing *model.Filing
	Docs   []*model.Document
}

var (
	txtURLRegex    = regexp.MustCompile(`(?s)<\s*td[^>]*>Complete submission text file<\s*/\s*td>.*<td[^>]*><a href="(.*?)">.*txt<\s*/\s*a><\s*/\s*td>`)
	relcikRegex    = regexp.MustCompile(`(?s)<span class="companyName">.*?\((.*?)...<acronym.*?CIK=(.*?)&`)
	altrelcikRegex = regexp.MustCompile(`(?s)<span class="companyName">.*?action=.*?">(.*?)</a>...<acronym.*?CIK=(.*?)&`)
	timeRegex      = regexp.MustCompile(`(?s)Accepted</div>.*?<div class="info">(.*?)</div>`)
)

func buildFiling(filer model.Filer, idxURL string) (filing ArchivedFiling, err error) {

	f := &model.Filing{
		Filer:      filer.Name,
		EdgarURL:   idxURL,
		CIK:        filer.CIK,
		Symbol:     filer.Symbol,
		AllCIKs:    []string{},
		AllSymbols: []string{},
	}
	// Need accession.
	f.Accession = getAccession(idxURL)

	indexPageHTML, err := getPage(idxURL, 2)
	if err != nil {
		return
	}
	// edgar time
	t, err := findTime(indexPageHTML)
	if err != nil {
		return
	}
	f.EdgarTime = t

	// relation
	relMap, err := findRelationshipMap(indexPageHTML)
	if err != nil {
		return
	}

	//
	rel := relMap[f.CIK]
	f.FilerRelation = rel

	for k := range relMap {
		f.AllCIKs = append(f.AllCIKs, k)
	}

	// form type
	docURL, err := findDocURL(indexPageHTML)
	if err != nil {
		return
	}

	docs, err := documents.GetDocsFromTxt("https://www.sec.gov" + docURL)
	if err != nil {
		return
	}

	f.FormType = docs[0].DocType

	return ArchivedFiling{Filing: f, Docs: docs}, nil
}

func getAccession(idxURL string) string {
	rs := strings.Split(idxURL, "/")
	idxString := rs[len(rs)-1]
	acc := strings.Split(idxString, "-index")
	return acc[0]
}

// findDocURL parses the text document url out of the index page.
func findDocURL(htmlstr string) (string, error) {
	urls := txtURLRegex.FindStringSubmatch(htmlstr)
	if len(urls) != 2 {
		return "", fmt.Errorf("did not parse doc url")
	}
	return urls[1], nil
}

// findRelationshipMap parses the rel-cik map out of the index page.
func findRelationshipMap(htmlstr string) (map[string]string, error) {
	matches := relcikRegex.FindAllStringSubmatch(htmlstr, -1)
	if matches == nil {
		return nil, fmt.Errorf("did not parse rel map")
	}

	if strings.Contains(matches[0][1], "<a") {
		matches = altrelcikRegex.FindAllStringSubmatch(htmlstr, -1)
		if matches == nil {
			return nil, fmt.Errorf("did not parse rel map")
		}
	}

	theMap := map[string]string{}

	for _, m := range matches {
		theMap[m[2]] = m[1]
	}

	return theMap, nil
}

// findRelationshipMap parses the rel-cik map out of the index page.
func findTime(htmlstr string) (time.Time, error) {
	matches := timeRegex.FindStringSubmatch(htmlstr)
	if matches == nil {
		return time.Time{}, fmt.Errorf("did not parse time")
	}
	return time.Parse("2006-01-02 15:04:05", matches[1])
}

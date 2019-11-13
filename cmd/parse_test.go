package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccesion(t *testing.T) {

	idxURL := "https://sec.gov/Archives/edgar/data/320193/000032019318000139/0000320193-18-000139-index.html"
	accession := getAccession(idxURL)
	assert.Equal(t, "0000320193-18-000139", accession)
}

func TestParseRelationships(t *testing.T) {
	file, err := os.Open("testdata/indexpage.htm")
	defer file.Close()
	assert.Nil(t, err, "could not open file")
	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err, "could not open file")
	text := string(data)

	//
	mp, err := findRelationshipMap(text)
	assert.Nil(t, err)
	assert.Equal(t, "Subject", mp["0000320193"])
}

func TestParseRelationshipsForm4(t *testing.T) {
	file, err := os.Open("testdata/indexpage-4.htm")
	defer file.Close()
	assert.Nil(t, err, "could not open file")
	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err, "could not open file")
	text := string(data)

	//
	mp, err := findRelationshipMap(text)
	assert.Nil(t, err)
	assert.Equal(t, "Issuer", mp["0000320193"])
}

func TestParseTime(t *testing.T) {
	file, err := os.Open("testdata/indexpage.htm")
	defer file.Close()
	assert.Nil(t, err, "could not open file")
	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err, "could not open file")
	text := string(data)

	//
	ttime, err := findTime(text)
	assert.Nil(t, err)
	assert.Equal(t, "1994-03-08 00:00:00 +0000 UTC", ttime.String())
}

func TestBuildFiling(t *testing.T) {
	//
	// filer := &model.Filer{
	// 	Name: "Apple, Inc",
	// 	CIK:  "0000320193",
	// }

	// idxURL := "https://www.sec.gov/Archives/edgar/data/320193/000119312513486406/0001193125-13-486406-index.html"

	// filing, err := buildFiling(filer, idxURL)
	// if assert.NoError(t, err) {
	// 	assert.Equal(t, filer.Name, filing.Filing.Filer)
	// 	fmt.Println(filing.Docs[0])
	// }
}

func TestExecute(t *testing.T) {
	// filer := &model.Filer{
	// 	Name: "Apple, Inc",
	// 	CIK:  "0000320193",
	// }
	// b := &Backfiller{}
	// b.execute(filer)
}

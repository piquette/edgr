package documents

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	ts := testServ()
	defer ts.Close()

	docs, err := Get(ts.URL + "/indexpage.htm")
	assert.Nil(t, err)
	assert.Len(t, docs, 2)
}

func TestParseText(t *testing.T) {

	file, err := os.Open("testdata/completefiling.txt")
	defer file.Close()
	assert.Nil(t, err, "could not open file")
	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err, "could not open file")
	text := string(data)

	//
	exampleURL := "https://example.com/"
	docs, err := parseFullText(exampleURL, text)
	assert.Nil(t, err)
	assert.Len(t, docs, 2)
	doc := docs[0]
	assert.Equal(t, "8-K", doc.DocType)
	assert.Equal(t, int64(1), doc.Sequence)
	assert.Equal(t, exampleURL+"a8k_20190122.htm", doc.EdgarURL)
	assert.Equal(t, "8-K", doc.Description)
	assert.Len(t, doc.Body, 28736)

	attchmnt := docs[1]
	assert.Equal(t, "EX-99.1", attchmnt.DocType)
	assert.Equal(t, int64(2), attchmnt.Sequence)
	assert.Equal(t, exampleURL+"a8k_ex991x20190122.htm", attchmnt.EdgarURL)
	assert.Equal(t, "EXHIBIT 99.1", attchmnt.Description)
	assert.Len(t, attchmnt.Body, 422636)
}

func TestParseForm4(t *testing.T) {
	ts := testServ()
	defer ts.Close()

	file, err := os.Open("testdata/form4.txt")
	defer file.Close()
	assert.Nil(t, err, "could not open file")
	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err, "could not open file")
	text := string(data)

	//
	docs, err := parseFullText(ts.URL+"/", text)
	assert.Nil(t, err)
	assert.Len(t, docs, 1)
	doc := docs[0]
	assert.Equal(t, "4", doc.DocType)
	assert.Equal(t, int64(1), doc.Sequence)
	assert.Equal(t, ts.URL+"/xslF345X03/edgar.xml", doc.EdgarURL)

	assert.Equal(t, "PRIMARY DOCUMENT", doc.Description)
}

func TestParseWithXRBL(t *testing.T) {
	file, err := os.Open("testdata/10q.txt")
	defer file.Close()
	assert.Nil(t, err, "could not open file")
	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err, "could not open file")
	text := string(data)

	//
	exampleURL := "https://example.com/"
	docs, err := parseFullText(exampleURL, text)
	assert.Nil(t, err)
	assert.Len(t, docs, 5)
}

func TestParseIndexPage(t *testing.T) {
	file, err := os.Open("testdata/indexpage.htm")
	defer file.Close()
	assert.Nil(t, err, "could not open file")
	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err, "could not open file")
	text := string(data)

	//
	url, err := findTextURL(text)
	assert.Nil(t, err)
	assert.Equal(t, "/Archives/edgar/data/1173431/000117343119000002/0001173431-19-000002.txt", url)
}

func TestParseIndexPageWithXBRL(t *testing.T) {
	file, err := os.Open("testdata/indexpage-with-xbrl.htm")
	defer file.Close()
	assert.Nil(t, err, "could not open file")
	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err, "could not open file")
	text := string(data)

	//
	url, err := findTextURL(text)
	assert.Nil(t, err)
	assert.Equal(t, "/Archives/edgar/data/318154/000031815419000008/0000318154-19-000008.txt", url)
}

func TestReplaceLinks(t *testing.T) {
	file, err := os.Open("testdata/document.html")
	defer file.Close()
	assert.Nil(t, err, "could not open file")
	data, err := ioutil.ReadAll(file)
	assert.Nil(t, err, "could not open file")
	text := string(data)

	url := "http://example.com/"
	replaced := replaceRelativeLinks(url, text)

	count := strings.Count(replaced, url)
	assert.Equal(t, 18, count)
}

// StartTestServer starts up a test server.
func testServ() *httptest.Server {
	path := `/testdata/`

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path = dir + path
	http.DefaultServeMux = new(http.ServeMux)
	fs := http.FileServer(http.Dir(path))
	return httptest.NewServer(fs)
}

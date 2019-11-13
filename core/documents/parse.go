package documents

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/piquette/edgr/core/model"
	"github.com/tdewolff/minify"
)

const (
	retryAttempts = 5
)

var (
	urlRegex       = regexp.MustCompile(`(?s)<\s*td[^>]*>Complete submission text file<\s*/\s*td>.*<td[^>]*><a href="(.*?)">.*txt<\s*/\s*a><\s*/\s*td>`)
	documentRegex  = regexp.MustCompile(`(?s)<\s*DOCUMENT[^>]*>(.*?)<\s*/\s*DOCUMENT>`)
	detailRegex    = regexp.MustCompile(`(?s)<TYPE>(.*?)<SEQUENCE>(.*?)<FILENAME>(.*?)<DESCRIPTION>(.*?)<TEXT>(.*?)<\s*/\s*TEXT>`)
	altDetailRegex = regexp.MustCompile(`(?s)<TYPE>(.*?)<SEQUENCE>(.*?)<FILENAME>(.*?)<TEXT>(.*?)<\s*/\s*TEXT>`)
	relativeRegex  = regexp.MustCompile(`\"[^\"]+\.(?i)(gif|png|jpg|jpeg)"`)
)

// Get gets the documents.
func Get(url string) (docs []*model.Document, err error) {
	indexPg, err := getPage(url, retryAttempts)
	if err != nil {
		return nil, err
	}

	textURL, err := findTextURL(indexPg)
	if err != nil {
		return nil, err
	}
	textURL = "https://www.sec.gov" + textURL

	return GetDocsFromTxt(textURL)
}

// GetDocsFromTxt gets documents without getting the index page first.
func GetDocsFromTxt(textURL string) (docs []*model.Document, err error) {

	textDocument, err := getPage(textURL, retryAttempts)
	if err != nil {
		return nil, err
	}

	//
	i := strings.LastIndex(textURL, "/")
	baseURL := textURL[0 : i+1]

	docs, err = parseFullText(baseURL, textDocument)
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("no docs found during parse")
	}

	for _, doc := range docs {
		// replace links.
		doc.Body = replaceRelativeLinks(baseURL, doc.Body)
		if doc.Body == "" {
			return nil, fmt.Errorf("doc body error during parse")
		}
	}
	return
}

// findTextURL parses the text document url out of the index page.
func findTextURL(html string) (string, error) {
	urls := urlRegex.FindStringSubmatch(html)
	if len(urls) != 2 {
		return "", fmt.Errorf("did not parse url")
	}
	return urls[1], nil
}

func parseFullText(url, fullText string) (docs []*model.Document, err error) {
	matches := documentRegex.FindAllStringSubmatch(fullText, -1)
	if matches == nil {
		return nil, fmt.Errorf("could not parse full text doc")
	}

	for _, match := range matches {
		if len(match) != 2 {
			log.Println("failed to parse a doc from txt file")
			continue
		}
		doctext := match[1]

		details := detailRegex.FindStringSubmatch(doctext)
		var desc string     // doc description.
		var bodytext string // the actual content.
		if len(details) != 6 {
			// Try alternate.
			details = altDetailRegex.FindStringSubmatch(doctext)
			if len(details) != 5 {
				log.Println("could not parse details from doc in txt file")
				continue
			} else {
				// Alternate success, but no description.
				bodytext = strings.TrimSpace(details[4])
				desc = "N/A"
			}
		} else {
			// Success.
			desc = strings.TrimSpace(details[4])
			bodytext = strings.TrimSpace(details[5])
		}

		docType := strings.TrimSpace(details[1])
		seq := strings.TrimSpace(details[2])
		filename := strings.TrimSpace(details[3])

		// Blacklist some files.
		if docType == "GRAPHIC" || strings.Contains(desc, "XBRL") {
			continue
		}

		edgarURL := url
		validfile, _ := regexp.MatchString(`.+(?i)\.(htm|html|txt)\z`, filename)
		if !validfile && strings.Contains(filename, ".xml") {
			// Weird file descriptor --
			var schema string
			switch docType {
			case "3", "3/A":
				schema = "xslF345X02/"
			case "4", "4/A":
				schema = "xslF345X03/"
			case "D", "D/A":
				schema = "xslFormDX01/"
			default:
				log.Println("weird doc type during replace:", docType)
				log.Println("doc url was:", edgarURL)
				continue
			}

			edgarURL = edgarURL + schema + filename
			htmlText, pageErr := getPage(edgarURL, 2)
			if pageErr != nil {
				log.Println("bad doc url during replace: ", edgarURL)
				continue
			}
			bodytext = htmlText

		} else if !validfile {
			log.Println("bad doc type: ", filename)
			continue
		} else {
			edgarURL = edgarURL + filename
		}

		// Minify.
		mini, err := minify.New().String("text/html", bodytext)
		if err == nil {
			bodytext = mini
		}

		// Create document.
		sequence, _ := strconv.Atoi(seq)
		doc := &model.Document{
			Sequence:    int64(sequence),
			DocType:     docType,
			Description: desc,
			EdgarURL:    edgarURL,
			Body:        bodytext,
		}

		docs = append(docs, doc)
	}
	return
}

func replaceRelativeLinks(url, text string) string {
	return relativeRegex.ReplaceAllStringFunc(text, func(match string) string {
		return `"` + url + match[1:]
	})
}

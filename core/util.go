package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// getPage executes a http GET and attempts retries.
func getPage(url string, retry int) (string, error) {
	for i := 0; i < retry; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			bytes, rerr := ioutil.ReadAll(resp.Body)
			return string(bytes), rerr
		}

		time.Sleep(1000 * time.Millisecond)
	}
	return "", fmt.Errorf("timed out retrieving %s", url)
}

// findListURLs parses the list of idx urls out of the directory page.
func findListURLs(html string) []string {
	matches := dirRegex.FindAllStringSubmatch(html, -1)
	if matches == nil || len(matches) == 1 {
		log.Warn("could not find matches")
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
		log.Warn("could not find matches")
		return ""
	}
	return "https://sec.gov" + matches[1] + "index.html"
}

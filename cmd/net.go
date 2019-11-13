package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

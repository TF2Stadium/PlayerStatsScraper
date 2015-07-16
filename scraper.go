package scraper

import (
	"net/http"

	"github.com/bitly/go-simplejson"
)

func getJsonFromUrl(url string) (*simplejson.Json, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := simplejson.NewFromReader(resp.Body)

	return data, err
}

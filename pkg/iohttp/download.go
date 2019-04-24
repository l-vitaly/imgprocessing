package iohttp

import (
	"io/ioutil"
	"net/http"
)

// GetContentByURL
func GetContentByURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

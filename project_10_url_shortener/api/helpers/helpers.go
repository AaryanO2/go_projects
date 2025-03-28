package helpers

import (
	"net/url"
	"os"
	"strings"
)

func RemoveDomainError(inputURL string) bool {
	paresedURL, err := url.Parse(inputURL)

	if err!= nil{
		return paresedURL.Host == os.Getenv("DOMAIN") 
	}

	newURL := strings.Replace(inputURL, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]
	return newURL != os.Getenv("DOMAIN")
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

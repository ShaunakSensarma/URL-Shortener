package helpers

import (
	"os"
	"strings"
)

// EnforceHTTP method ensures that all request are going over HTTP.
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

/*
this function is likely intended to determine whether a given URL should have its domain removed
based on a comparison with a predefined domain stored in the environment variable.

We don't want the user to access localhost:3000.
*/
func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	// If the input URL doesn't match exactly, remove common URL prefixes (http://, https://, www://) from the URL.
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www://", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	return newURL != os.Getenv("DOMAIN")
}

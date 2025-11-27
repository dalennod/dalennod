package archive

import (
	"log"
	"net/http"

	"dalennod/internal/default_client"
)

func SendSnapshot(url string) (bool, string) {
	var check bool = checkURL(url)
	var snapshotURL string

	if check {
		urlToRequest := "https://web.archive.org/save/" + url
		res, err := default_client.HttpDefaultClientDo(http.MethodGet, urlToRequest)
		if err != nil {
			log.Println("WARN: Failed to archive due to error: ", err)
			return false, "Failed to archive due to error"
		}
		defer res.Body.Close()

		snapshotURL = res.Request.URL.String()
		log.Printf("INFO: Archived URL '%s' to Wayback Machine: %s\n", url, snapshotURL)

		return check, snapshotURL
	} else {
		log.Printf("WARN: URL '%s' did not respond. Not sending to be archived.\n", url)
		return check, "Failed to archive due to website not responding"
	}
}

func checkURL(url string) bool {
	res, err := default_client.HttpDefaultClientDo(http.MethodGet, url)
	if err != nil {
		log.Println("WARN: Failed to ping website", err)
		return false
	}
	defer res.Body.Close()

	return true
}

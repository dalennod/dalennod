package archive

import (
	"dalennod/internal/logger"
	"fmt"
	"net/http"
)

func SendSnapshot(url string) (bool, string) {
	var check bool = checkURL(url)
	var snapshotURL string

	if check {
		res, err := http.Get(fmt.Sprintf("https://web.archive.org/save/%s", url))
		if err != nil {
			logger.Warn.Println("Failed to archive due to error: ", err)
			return false, "Failed to archive due to error"
		}
		defer res.Body.Close()

		snapshotURL = res.Request.URL.String()
		logger.Info.Printf("Archived URL '%s' to Wayback Machine. [%s]\n", url, snapshotURL)

		return check, snapshotURL
	} else {
		logger.Warn.Printf("URL [%s] did not respond. Not sending to be archived.\n", url)
		return check, "Failed to archive due to website not responding"
	}
}

func checkURL(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		logger.Warn.Println("Failed to ping website", err)
	}

	if res != nil {
		defer res.Body.Close()
		return true
	} else {
		return false
	}
}

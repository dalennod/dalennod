package archive

import (
	"dalennod/internal/logger"
	"fmt"
	"net/http"
)

func SendSnapshot(url string) {
	var check bool = checkURL(url)
	if check {
		res, err := http.Get(fmt.Sprintf("https://web.archive.org/save/%s", url))
		if err != nil {
			logger.Warn.Println("Failed to archive due to error: ", err)
		}
		defer res.Body.Close()

		logger.Info.Printf("Archived URL '%s' to Wayback Machine. [%s]\n", url, res.Request.URL.String())
	} else {
		logger.Warn.Printf("URL [%s] did not respond. Not sending to be archived.", url)
	}
}

func checkURL(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		logger.Warn.Println("Failed to ping website.", err)
	}

	if res != nil {
		defer res.Body.Close()
		return true
	} else {
		return false
	}
}

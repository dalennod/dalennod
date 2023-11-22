package archive

import (
	"dalennod/internal/logger"
	"fmt"
	"net/http"
)

func SendSnapshot(url string) {
	res, err := http.Get(fmt.Sprintf("https://web.archive.org/save/%s", url))
	if err != nil {
		logger.Warn.Println("Failed to archive due to error: ", err)
	}
	defer res.Body.Close()

	logger.Info.Printf("Archived URL '%s' to Wayback Machine. [%s]\n", url, res.Request.URL.String())
}

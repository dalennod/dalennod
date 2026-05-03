package server

import (
	"fmt"
	"html/template"
	"regexp"

	"dalennod/internal/constants"
)

func getHostname(input string) string {
	hostnamePattern := regexp.MustCompile(`^https?:\/\/([^\/]+)`)
	matches := hostnamePattern.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func pageTitle() string {
	return webPageTitle
}

func keywordSplit(keywords string, delimiter string) []string {
	re := regexp.MustCompile(`\s*,\s*`) // To accomplish whitespace trimming without additional loops
	return re.Split(keywords, -1)
}

func grabThumbnail(id int64) string {
	return fmt.Sprintf("http://localhost%s/thumbnail/%d", constants.SECONDARY_PORT, id)
}

func encapsulateURL(bkmURL string) template.URL {
	return template.URL(bkmURL)
}

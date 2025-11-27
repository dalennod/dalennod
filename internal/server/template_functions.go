package server

import (
	"fmt"
	"html/template"
	"regexp"

	"dalennod/internal/constants"
)

func getHostname(input string) string {
	hostnamePattern := regexp.MustCompile(`^(?:http(?:s?):\/\/(?:www\.)?)?([A-Za-z0-9_:.-]+)\/?`)
	matches := hostnamePattern.FindAllString(input, -1)
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}

func webUIAddress() string {
	if constants.WEBUI_ADDR[0] == 58 { // ':'
		return "http://localhost" + constants.WEBUI_ADDR
	} else {
		return "http://%s\n" + constants.WEBUI_ADDR
	}
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

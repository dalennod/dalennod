package server

import (
    "fmt"
    "regexp"

    "dalennod/internal/constants"
)

func getHostname(input string) string {
    // hostnamePattern := regexp.MustCompile(`(?i)(?:(?:https?|ftp):\/\/)?(?:www\.)?(?:[a-z0-9]([-a-z0-9]*[a-z0-9])?\.)+[a-z]{2,63}`)
    hostnamePattern := regexp.MustCompile(`^(?:http(?:s?):\/\/(?:www\.)?)?([A-Za-z0-9_:.-]+)\/?`) // Will match localhost and ports too
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
    // return strings.Split(keywords, delimiter)
    re := regexp.MustCompile(`\s*,\s*`) // To accomplish whitespace trimming without additional loops
    return re.Split(keywords, -1)
}

func grabThumbnail(id int) string {
    return fmt.Sprintf("http://localhost%s/thumbnail/%d", constants.SECONDARY_PORT, id)
}

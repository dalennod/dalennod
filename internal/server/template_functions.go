package server

import (
    "dalennod/internal/constants"
    "encoding/base64"
    "net/http"
    "regexp"
    "strings"
)

func getHostname(input string) string {
    hostnamePattern := regexp.MustCompile(`(?i)(?:(?:https?|ftp):\/\/)?(?:www\.)?(?:[a-z0-9]([-a-z0-9]*[a-z0-9])?\.)+[a-z]{2,63}`)
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
    return strings.Split(keywords, delimiter)
}

func byteConversion(blobImage []byte) string {
    var base64Encoded string

    mimeType := http.DetectContentType(blobImage)
    switch mimeType {
    case "image/avif":
        base64Encoded = "avif;base64,"
    case "image/webp":
        base64Encoded = "webp;base64,"
    case "image/png":
        base64Encoded = "png;base64,"
    case "image/jpeg":
        base64Encoded = "jpeg;base64,"
    default:
        base64Encoded = "jpeg;base64,"
    }
    base64Encoded += base64.StdEncoding.EncodeToString(blobImage)

    return base64Encoded
}

func pageCountUp() int {
    pageCount = pageCount + 2
    return pageCount
}

func pageCountDown() int {
    pageCount = pageCount - 1
    return pageCount
}

func pageCountNowUpdate() int {
    return pageCount - 1
}

func pageCountNowDelete() int {
    return pageCount
}

package thumb_url

import (
	"io"
	"net/http"
	"strings"

	"github.com/dyatlov/go-opengraph/opengraph"
)

// can be improved in future to use overall less code by using checkURL from archive package
func GetPageThumb(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	pageHtml, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var og *opengraph.OpenGraph = opengraph.NewOpenGraph()
	err = og.ProcessHTML(strings.NewReader(string(pageHtml)))
	if err != nil {
		return "", err
	}

	if len(og.Images) == 0 {
		return "", nil
	}

	var thumbURL string = og.Images[0].URL

	return thumbURL, nil
}

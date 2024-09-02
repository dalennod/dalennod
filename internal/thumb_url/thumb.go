// This code was inspired by and is a simplified version of Vitali Deatlov's opengraph package.
// pkg: https://pkg.go.dev/github.com/dyatlov/go-opengraph/opengraph

package thumb_url

import (
	"dalennod/internal/logger"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type OGData struct {
	URL              string    `json:"url"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Determiner       string    `json:"determiner"`
	SiteName         string    `json:"site_name"`
	Locale           string    `json:"locale"`
	LocalesAlternate []string  `json:"locales_alternate"`
	Images           []*Images `json:"images"`
}

type Images struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"`
}

func GetPageThumb(url string) (string, []byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()

	pageHtml, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil, err
	}

	var og *OGData = &OGData{}
	err = og.readHTML(strings.NewReader(string(pageHtml)))
	if err != nil {
		return "", nil, err
	}

	if len(og.Images) == 0 {
		return "", nil, nil
	}

	var thumbURL string = og.Images[0].URL

	if thumbURL == "" {
		return thumbURL, nil, nil
	} else {
		var byteThumbURL []byte = getBase64(thumbURL)
		return thumbURL, byteThumbURL, nil
	}
}

func getBase64(thumbURL string) []byte {
	var thumbUrlBytes []byte
	resp, err := http.Get(thumbURL)
	if err != nil {
		logger.Warn.Println("could not request thumburl")
		return thumbUrlBytes
	}
	defer resp.Body.Close()

	thumbUrlBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn.Println("could not read thumburl")
	}

	return thumbUrlBytes
}

func (og *OGData) readHTML(buffer io.Reader) error {
	var htmlToken *html.Tokenizer = html.NewTokenizer(buffer)
	for {
		var tokenType html.TokenType = htmlToken.Next()
		switch tokenType {
		case html.ErrorToken:
			if htmlToken.Err() == io.EOF {
				return nil
			}
			return htmlToken.Err()
		case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken:
			name, hasAttribute := htmlToken.TagName()
			if atom.Lookup(name) != atom.Meta || !hasAttribute {
				continue
			}
			var valMap map[string]string = make(map[string]string)
			var key, val []byte
			for hasAttribute {
				key, val, hasAttribute = htmlToken.TagAttr()
				valMap[atom.String(key)] = string(val)
			}
			og.readMeta(valMap)
		}
	}
}

func (og *OGData) readMeta(metaAttributes map[string]string) {
	switch metaAttributes["property"] {
	case "og:description":
		og.Description = metaAttributes["content"]
	case "og:title":
		og.Title = metaAttributes["content"]
	case "og:url":
		og.URL = metaAttributes["content"]
	case "og:determiner":
		og.Determiner = metaAttributes["content"]
	case "og:site_name":
		og.SiteName = metaAttributes["content"]
	case "og:locale":
		og.Locale = metaAttributes["content"]
	case "og:locale:alternate":
		og.LocalesAlternate = append(og.LocalesAlternate, metaAttributes["content"])
	case "og:image":
		og.Images = addImageUrl(og.Images, metaAttributes["content"])
	case "og:image:url":
		og.Images = addImageUrl(og.Images, metaAttributes["content"])
	case "og:image:secure_url":
		og.Images = addImageSecureUrl(og.Images, metaAttributes["content"])
	case "og:image:type":
		og.Images = addImageType(og.Images, metaAttributes["content"])
	default:
		return
	}
}

func addImageUrl(images []*Images, v string) []*Images {
	if len(images) == 0 || (images[len(images)-1].URL != "" && images[len(images)-1].URL != v) {
		images = append(images, &Images{})
	}
	images[len(images)-1].URL = v
	return images
}

func addImageSecureUrl(images []*Images, v string) []*Images {
	images = ensureHasImage(images)
	images[len(images)-1].SecureURL = v
	return images
}

func ensureHasImage(images []*Images) []*Images {
	if len(images) == 0 {
		images = append(images, &Images{})
	}
	return images
}

func addImageType(images []*Images, v string) []*Images {
	images = ensureHasImage(images)
	images[len(images)-1].Type = v
	return images
}

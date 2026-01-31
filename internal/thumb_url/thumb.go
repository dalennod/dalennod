// This code was inspired by and is a simplified version of Vitali Deatlov's opengraph package. Thank you
// pkg: https://pkg.go.dev/github.com/dyatlov/go-opengraph/opengraph
//
// This code was inspired by and is a version of Adam Presley's Go Favicon Grabber. Thank you
// repo: https://github.com/adampresley/gofavigrab

package thumb_url

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"dalennod/internal/constants"
	"dalennod/internal/default_client"
)

type OGData struct {
	// URL              string    `json:"url"`
	// Title            string    `json:"title"`
	// Description      string    `json:"description"`
	// Determiner       string    `json:"determiner"`
	// SiteName         string    `json:"site_name"`
	// Locale           string    `json:"locale"`
	// LocalesAlternate []string  `json:"locales_alternate"`
	Images []*Images `json:"images"`
}

type Images struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"`
}

func GetPageThumb(url string) (string, error) {
	res, err := default_client.HttpDefaultClientDo(http.MethodGet, url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	pageHtml, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var og *OGData = &OGData{}
	err = og.readHTML(strings.NewReader(string(pageHtml)))
	if err != nil {
		return "", err
	}

	if len(og.Images) > 0 && og.Images[0].URL != "" {
		thumbURL := og.Images[0].URL
		return thumbURL, nil
	}

	rawFaviconURL, err := getPageFavicon(string(pageHtml))
	if err != nil {
		return "", err
	}

	resolvedFaviconURL, err := normalizeFaviconURL(url, rawFaviconURL)
	if err != nil {
		return "", fmt.Errorf("did not find any thumbnail or favicon in webpage")
	}

	return resolvedFaviconURL, nil
}

func normalizeFaviconURL(baseURL, rawURL string) (string, error) {
	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	resolvedURL := parsedBase.ResolveReference(parsedURL)
	result := resolvedURL.String()

	if strings.Contains(result, ".svg") {
		return "", fmt.Errorf("SVG favicon is not supported")
	}

	return result, nil
}

func getPageFavicon(url string) (string, error) {
	tokenizer := html.NewTokenizer(strings.NewReader(url))

	var tokenType html.TokenType
	var hasAttributes bool

	var attributeKey []byte
	var attributeValue []byte
	var hasMoreAttributes bool

	hasFavicon := false

	for {
		tokenType = tokenizer.Next()
		if tokenType == html.ErrorToken {
			log.Println("WARN: error while parsing HTML:", tokenizer.Err())
			break
		}
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			_, hasAttributes = tokenizer.TagName()
			if hasAttributes {
				for {
					attributeKey, attributeValue, hasMoreAttributes = tokenizer.TagAttr()
					if string(attributeKey) == "ref" || string(attributeKey) == "rel" {
						if strings.Contains(string(attributeValue), "shortcut") || strings.Contains(string(attributeValue), "icon") {
							hasFavicon = true
						}
					}
					if string(attributeKey) == "href" && hasFavicon {
						// This only returns the raw value. Need to account for non-normalized URLs
						return string(attributeValue), nil
					}
					if !hasMoreAttributes {
						break
					}
				}
			}
		}
	}

	return "", fmt.Errorf("URL not found")
}

func adjustDownloadedThumbnailSize(sourceFilePath string, width int) error {
	inputSource, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer inputSource.Close()

	bytes, err := io.ReadAll(inputSource)
	if err != nil {
		return err
	}
	mimetype := http.DetectContentType(bytes)
	inputSource.Seek(0, 0)

	var source image.Image

	switch mimetype {
	case "image/png":
		source, err = png.Decode(inputSource)
	case "image/jpeg":
		source, err = jpeg.Decode(inputSource)
	default:
		log.Printf("INFO: Will not resize because the image type is %s: which is unsupported\n", mimetype)
		return nil
	}
	if err != nil {
		return err
	}

	sourceBounds := source.Bounds()
	sourceW := sourceBounds.Dx()
	if sourceW <= width {
		log.Printf("INFO: Image size already at or below width: %d\n", width)
		return nil
	}

	ratio := (float64)(sourceBounds.Max.Y) / (float64)(sourceBounds.Max.X)
	height := int(math.Round(float64(width) * ratio))
	dest := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dest, dest.Rect, source, sourceBounds, draw.Over, nil)

	resizedImagePath := sourceFilePath+"_resized"
	resizedImage, err := os.Create(resizedImagePath)
	if err != nil {
		return err
	}

	switch mimetype {
	case "image/png":
		err = png.Encode(resizedImage, dest)
	case "image/jpeg":
		err = jpeg.Encode(resizedImage, dest, nil)
	}

	if err != nil {
		os.Remove(resizedImagePath)
		return err
	}

	os.Remove(sourceFilePath)
	os.Rename(resizedImagePath, sourceFilePath)

	return nil
}

func DownThumb(id int64, thumbURL string) error {
	bookmarkIDStr := strconv.FormatInt(id, 10)
	outputFilePath := filepath.Join(constants.THUMBNAILS_PATH, bookmarkIDStr)
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	fileRequest, err := default_client.HttpDefaultClientDo(http.MethodGet, thumbURL)
	if err != nil {
		return err
	}
	defer fileRequest.Body.Close()

	_, err = io.Copy(outFile, fileRequest.Body)
	if err != nil {
		return err
	}

	if err := adjustDownloadedThumbnailSize(outputFilePath, 300); err != nil {
		return err
	}

	return nil
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
	// case "og:description":
	//     og.Description = metaAttributes["content"]
	// case "og:title":
	//     og.Title = metaAttributes["content"]
	// case "og:url":
	//     og.URL = metaAttributes["content"]
	// case "og:determiner":
	//     og.Determiner = metaAttributes["content"]
	// case "og:site_name":
	//     og.SiteName = metaAttributes["content"]
	// case "og:locale":
	//     og.Locale = metaAttributes["content"]
	// case "og:locale:alternate":
	//     og.LocalesAlternate = append(og.LocalesAlternate, metaAttributes["content"])
	case "og:image":
		og.Images = addImageUrl(og.Images, metaAttributes["content"])
	case "og:image:url":
		og.Images = addImageUrl(og.Images, metaAttributes["content"])
	case "og:image:secure_url":
		og.Images = addImageSecureUrl(og.Images, metaAttributes["content"])
	// case "og:image:type":
	// 	og.Images = addImageType(og.Images, metaAttributes["content"])
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

// func addImageType(images []*Images, v string) []*Images {
// 	images = ensureHasImage(images)
// 	images[len(images)-1].Type = v
// 	return images
// }

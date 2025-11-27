package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"dalennod/internal/bookmark_import"
	"dalennod/internal/constants"
	"dalennod/internal/db"
)

func importHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl = template.Must(template.New("import").ParseFS(Web, "web/import/index.html"))
		if err := tmpl.ExecuteTemplate(w, "import", nil); err != nil {
			log.Println("wARN: failed to execute import template:", err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func importBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(constants.IMPORT_FILE_SIZE); err != nil {
			log.Println("ERROR: parsing form while importing bookmark:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		importFile, _, err := r.FormFile("importFile")
		if err != nil {
			log.Println("ERROR: parsing import file:", err)
			w.Write([]byte("error parsing import file. ERROR: " + err.Error()))
			return
		}
		defer importFile.Close()

		selectedBrowser := r.FormValue("selectedBrowser")

		if selectedBrowser == "Firefox" {
			firefoxBookmarks := &bookmark_import.Item{}
			if err := json.NewDecoder(importFile).Decode(firefoxBookmarks); err != nil {
				log.Println("WARN: could not decode imported bookmarks:", err)
				return
			}

			parsedBookmarks := bookmark_import.ParseFirefox(firefoxBookmarks, "")
			parsedBookmarksLength := strconv.Itoa(len(parsedBookmarks))

			for _, parsedBookmark := range parsedBookmarks {
				db.Add(database, parsedBookmark)
				output := "Added || { TITLE: " + parsedBookmark.Title + ", URL: " + parsedBookmark.URL + "CATEGORY: " + parsedBookmark.Category + ", KEYWORDS: " + parsedBookmark.Keywords + "}\n"
				w.Write([]byte(output))
			}
			w.Write([]byte("Added " + parsedBookmarksLength + " bookmarks to database."))
			return
		} else if selectedBrowser == "Chromium" {
			var chromiumBookmarks bookmark_import.ChromiumItem
			if err := json.NewDecoder(importFile).Decode(&chromiumBookmarks); err != nil {
				log.Println("WARN: could not decode imported bookmarks:", err)
				return
			}

			parsedBookmarks := bookmark_import.ParseChromium(chromiumBookmarks)
			parsedBookmarksLength := strconv.Itoa(len(parsedBookmarks))

			for _, parsedBookmark := range parsedBookmarks {
				db.Add(database, parsedBookmark)
				output := "Added || { TITLE: " + parsedBookmark.Title + ", URL: " + parsedBookmark.URL + "CATEGORY: " + parsedBookmark.Category + ", KEYWORDS: " + parsedBookmark.Keywords + "}\n"
				w.Write([]byte(output))
			}
			w.Write([]byte("Added " + parsedBookmarksLength + " bookmarks to database."))
			return
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

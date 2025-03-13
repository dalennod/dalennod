package server

import (
    "dalennod/internal/bookmark_import"
    "dalennod/internal/db"
    "dalennod/internal/logger"
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "strconv"
)

func importHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl = template.Must(template.New("import").ParseFS(Web, "web/import/index.html"))
        if err := tmpl.ExecuteTemplate(w, "import", nil); err != nil {
            logger.Warn.Println(err)
        }
    } else {
        internalServerErrorHandler(w, r)
    }
}

func importBookmarkHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 << 20 = 10 * (2^20) = 10.485.760 = ~10,48MB file size limit
            logger.Error.Println("error parsing form while importing bookmark. ERROR:", err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        importFile, _, err := r.FormFile("importFile")
        if err != nil {
            logger.Error.Println("error parsing import file. ERROR:", err)
            w.Write([]byte("error parsing import file. ERROR: " + err.Error()))
            return
        }
        defer importFile.Close()

        selectedBrowser := r.FormValue("selectedBrowser")

        if selectedBrowser == "Firefox" {
            firefoxBookmarks := &bookmark_import.Item{}
            if err := json.NewDecoder(importFile).Decode(firefoxBookmarks); err != nil {
                logger.Error.Println(err)
                fmt.Println(err)
                return
            }

            parsedBookmarks := bookmark_import.ParseFirefox(firefoxBookmarks, "")
            parsedBookmarksLength := strconv.Itoa(len(parsedBookmarks))

            for _, parsedBookmark := range parsedBookmarks {
                db.Add(database, parsedBookmark)
                output := "Added || { TITLE: " + parsedBookmark.Title + ", URL: " + parsedBookmark.URL + "GROUP: " + parsedBookmark.BmGroup + ", KEYWORDS: " + parsedBookmark.Keywords + "}\n"
                w.Write([]byte(output))
            }
            w.Write([]byte("Added " + parsedBookmarksLength + " bookmarks to database."))
            return
        } else if selectedBrowser == "Chromium" {
            var chromiumBookmarks bookmark_import.ChromiumItem
            if err := json.NewDecoder(importFile).Decode(&chromiumBookmarks); err != nil {
                logger.Error.Println(err)
                fmt.Println(err)
                return
            }

            parsedBookmarks := bookmark_import.ParseChromium(chromiumBookmarks)
            parsedBookmarksLength := strconv.Itoa(len(parsedBookmarks))

            for _, parsedBookmark := range parsedBookmarks {
                db.Add(database, parsedBookmark)
                output := "Added || { TITLE: " + parsedBookmark.Title + ", URL: " + parsedBookmark.URL + "GROUP: " + parsedBookmark.BmGroup + ", KEYWORDS: " + parsedBookmark.Keywords + "}\n"
                w.Write([]byte(output))
            }
            w.Write([]byte("Added " + parsedBookmarksLength + " bookmarks to database."))
            return
        }
    } else {
        internalServerErrorHandler(w, r)
    }
}

package server

import (
	"dalennod/internal/archive"
	"dalennod/internal/bookmark_import"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"
)

const PORT string = ":41415"

var (
	pageCount    int                    = 0
	tmplFuncMap  template.FuncMap       = make(template.FuncMap)
	allBookmarks map[string]interface{} = make(map[string]interface{})
	database     *sql.DB
	tmpl         *template.Template
	Web          embed.FS
)

func Start(data *sql.DB) {
	database = data

	var mux *http.ServeMux = http.NewServeMux()

	var fsopen fs.FS = fs.FS(Web)
	webStatic, err := fs.Sub(fsopen, "web/static")
	if err != nil {
		logger.Error.Fatalln(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(webStatic))))

	tmplFuncMap["getHostname"] = getHostname
	tmplFuncMap["keywordSplit"] = keywordSplit
	tmplFuncMap["byteConversion"] = byteConversion
	tmplFuncMap["pageCountUp"] = pageCountUp
	tmplFuncMap["pageCountDown"] = pageCountDown
	tmplFuncMap["pageCountNowUpdate"] = pageCountNowUpdate
	tmplFuncMap["pageCountNowDelete"] = pageCountNowDelete

	mux.HandleFunc("/{$}", rootHandler)
	mux.HandleFunc("/import/", importHandler)
	mux.HandleFunc("/api/import-bookmark/", importBookmarkHandler)
	mux.HandleFunc("/api/add/", addHandler)
	mux.HandleFunc("/api/row/{id}", rowHandler)
	mux.HandleFunc("/api/update/{id}", updateHandler)
	mux.HandleFunc("/api/delete/{id}", deleteHandler)
	mux.HandleFunc("/api/groups/", groupsHandler)
	mux.HandleFunc("/api/search/", searchHandler)
	mux.HandleFunc("/api/search-keyword/", searchKeywordHandler)
	mux.HandleFunc("/api/search-group/", searchGroupHandler)
	mux.HandleFunc("/api/search-hostname/", searchHostnameHandler)
	mux.HandleFunc("/api/check-url/", checkUrlHandler)
	mux.HandleFunc("/api/refetch-thumbnail/{id}", refetchThumbnailHandler)

	logger.Info.Printf("Web-server starting on http://localhost%s/\n", PORT)
	fmt.Printf("Web-server starting on http://localhost%s/\n", PORT)

	if err := http.ListenAndServe(PORT, mux); err != nil {
		fmt.Printf("Stopping (error: %v)\n", err)
		logger.Error.Printf("Stopping (error: %v)\n", err)
	}
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
	logger.Warn.Printf("status 500 at '%s%s'\n", r.Host, r.URL)
}

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
			jsonDecoder := json.NewDecoder(importFile)
			jsonDecoder.Decode(firefoxBookmarks)

			parsedBookmarks := bookmark_import.ParseFirefox(firefoxBookmarks, "")
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		pageCount = 0
		var bookmarks []setup.Bookmark
		r.ParseForm()

		var pageNo string = r.FormValue("page")
		if pageNo == "" {
			bookmarks = db.ViewAllWebUI(database, 0)
		} else {
			pageNoInt, err := strconv.Atoi(pageNo)
			if err != nil {
				logger.Error.Printf("error: invalid page no. %v", err)
			}
			pageCount = pageNoInt
			bookmarks = db.ViewAllWebUI(database, pageNoInt)
		}

		tmpl = template.Must(template.New("index").Funcs(tmplFuncMap).ParseFS(Web, "web/index.html"))
		allBookmarks["Bookmarks"] = bookmarks
		if err := tmpl.ExecuteTemplate(w, "index", allBookmarks); err != nil {
			logger.Warn.Println("error executing template for root index. ERROR:", err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodGet {
		matchId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			return
		}

		db.Remove(database, matchId)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodPost {
		var insData setup.Bookmark
		if err := json.NewDecoder(r.Body).Decode(&insData); err != nil {
			logger.Error.Println(err)
		}

		if !insData.Archived {
			db.Add(database, insData)
		} else {
			insData.Archived, insData.SnapshotURL = archive.SendSnapshot(insData.URL)
			if insData.Archived {
				db.Add(database, insData)
			} else {
				logger.Warn.Println("Snapshot failed")
				db.Add(database, insData)
			}
		}

		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodGet {
		w.Write([]byte("Alive"))
	} else {
		internalServerErrorHandler(w, r)
	}
}

func rowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		matchId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			return
		}

		oldData, err := db.ViewSingleRow(database, matchId, true)
		if err != nil {
			logger.Error.Println(err)
		}

		writeData, err := json.Marshal(&oldData)
		if err != nil {
			logger.Error.Println(err)
			return
		}

		w.Write([]byte(writeData))
	} else {
		internalServerErrorHandler(w, r)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodPost {
		matchId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			return
		}

		var newData setup.Bookmark
		err = json.NewDecoder(r.Body).Decode(&newData)
		if err != nil {
			logger.Error.Println(err)
		}
		newData.ID = matchId

		if !newData.Archived {
			db.Update(database, newData, true)
		} else {
			newData.Archived, newData.SnapshotURL = archive.SendSnapshot(newData.URL)
			if newData.Archived {
				db.Update(database, newData, true)
			} else {
				logger.Warn.Println("Snapshot failed")
				db.Update(database, newData, true)
			}
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func findSearchTerm(r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}
	searchTerm := r.FormValue("searchTerm")
	return searchTerm, nil
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodPost {
		searchTerm, err := findSearchTerm(r)
		if err != nil {
			logger.Error.Println("error getting search term. ERROR:", err)
		}

		if searchTerm == "" {
			rootHandler(w, &http.Request{
				Method: http.MethodGet,
				Form:   make(url.Values),
			})
			return
		}

		var bookmarks []setup.Bookmark = db.ViewAllWhere(database, searchTerm)
		allBookmarks["Bookmarks"] = bookmarks
		if err := tmpl.ExecuteTemplate(w, "bm_list", allBookmarks); err != nil {
			logger.Error.Println(err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func searchKeywordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodPost {
		searchTerm, err := findSearchTerm(r)
		if err != nil {
			logger.Error.Println("error getting search term. ERROR:", err)
		}

		var bookmarks []setup.Bookmark = db.ViewAllWhereKeyword(database, searchTerm)
		allBookmarks["Bookmarks"] = bookmarks
		if err := tmpl.ExecuteTemplate(w, "bm_list", allBookmarks); err != nil {
			logger.Error.Println(err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func searchGroupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodPost {
		searchTerm, err := findSearchTerm(r)
		if err != nil {
			logger.Error.Println("error getting search term. ERROR:", err)
		}

		var bookmarks []setup.Bookmark = db.ViewAllWhereGroup(database, searchTerm)
		allBookmarks["Bookmarks"] = bookmarks
		if err := tmpl.ExecuteTemplate(w, "bm_list", allBookmarks); err != nil {
			logger.Error.Println(err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func searchHostnameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		searchTerm, err := findSearchTerm(r)
		if err != nil {
			logger.Error.Println("error getting search term. ERROR:", err)
		}

		bookmarks := db.ViewAllWhereHostname(database, searchTerm)
		allBookmarks["Bookmarks"] = bookmarks
		if err := tmpl.ExecuteTemplate(w, "bm_list", allBookmarks); err != nil {
			logger.Error.Printf("error executing template after searching hostname. ERROR:%v\n", err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func checkUrlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodPost {
		var getData setup.Bookmark
		var err error = json.NewDecoder(r.Body).Decode(&getData)
		if err != nil {
			logger.Error.Println(err)
		}

		var searchUrl string = getData.URL
		getData, err = db.SearchByUrl(database, searchUrl)
		if err != nil {
			logger.Error.Println(err)
		}

		if getData.ID == 0 {
			http.NotFound(w, r)
			return
		}

		writeData, err := json.Marshal(&getData)
		if err != nil {
			logger.Error.Println(err)
		}
		w.Write([]byte(writeData))
	} else {
		internalServerErrorHandler(w, r)
	}
}

func refetchThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodGet {
		matchId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			return
		}

		if err := db.RefetchThumbnail(database, matchId); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Thumbnail updated"))
	} else {
		internalServerErrorHandler(w, r)
	}
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodGet {
		var currGroups map[string]interface{} = make(map[string]interface{})
		listCurrGroups, err := db.GetAllGroups(database)
		if err != nil {
			logger.Error.Println(err)
		}
		currGroups["AllGroups"] = listCurrGroups

		for _, group := range listCurrGroups {
			fmt.Fprintf(w, "<option value=\"%s\"></option>", group)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

package server

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"

	"dalennod/internal/archive"
	"dalennod/internal/constants"
	"dalennod/internal/db"
	"dalennod/internal/setup"
)

var (
	pageCountForSearch int                         = 0
	tmplFuncMap        template.FuncMap            = make(template.FuncMap)
	bookmarksMap       map[string][]setup.Bookmark = make(map[string][]setup.Bookmark)
	database           *sql.DB
	tmpl               *template.Template
	webPageTitle       string
	Web                embed.FS
)

type GotURLParams struct {
	pageNumber string
	searchType string
	searchTerm string
}

func parseURLParams(r *http.Request) GotURLParams {
	var gotURLParams GotURLParams
	err := r.ParseForm()
	if err != nil {
		log.Println("WARN: error parsing URL params:", err)
		return gotURLParams
	}
	gotURLParams.pageNumber = r.FormValue("page")
	gotURLParams.searchType = r.FormValue("search-type")
	gotURLParams.searchTerm = r.FormValue("search-term")
	return gotURLParams
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		pageCount := 0
		var bookmarks []setup.Bookmark
		var recentlyInteracted setup.RecentInteractions

		webPageTitle = "Dalennod"

		gotURLParams := parseURLParams(r)

		if gotURLParams.pageNumber != "" {
			webPageTitle = fmt.Sprintf("Page %s | Dalennod", gotURLParams.pageNumber)
			pageNoValid, err := strconv.Atoi(gotURLParams.pageNumber)
			if err != nil {
				log.Println("WARN: Invalid page number. Got:", gotURLParams.pageNumber)
			}
			pageCount = pageNoValid
		}

		if gotURLParams.searchType != "" {
			webPageTitle = "Search | Dalennod"
			switch gotURLParams.searchType {
			case "general":
				openPrefix := "o "
				searchTermAfter, prefixFound := strings.CutPrefix(gotURLParams.searchTerm, openPrefix)
				if prefixFound && searchTermAfter != "" {
					gotURLParams.searchTerm = searchTermAfter
					bookmark := db.OpenSesame(database, gotURLParams.searchTerm)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(bookmark)
					return
				}
				fallthrough
			case "hostname":
				fallthrough
			case "keyword":
				fallthrough
			case "category":
				bookmarks, pageCountForSearch = db.SearchFor(database, gotURLParams.searchType, gotURLParams.searchTerm, pageCount)
			default:
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "ERROR: Unrecognized search type")
				log.Println("ERROR: Unrecognized search type. Got:", gotURLParams.searchType)
				return
			}
		} else {
			bookmarks = db.ViewAllWebUI(database, pageCount)
		}

		if gotURLParams.pageNumber == "" && gotURLParams.searchType == "" {
			recentlyInteracted = db.RecentInteractions(database)
		}

		tmpl = template.Must(template.New("index").Funcs(tmplFuncMap).ParseFS(Web, "web/index.html"))
		bookmarksMap["AllBookmarks"] = bookmarks
		bookmarksMap["RecentBookmarks"] = recentlyInteracted.Bookmarks
		if err := tmpl.ExecuteTemplate(w, "index", bookmarksMap); err != nil {
			log.Println("ERROR: executing template for root index:", err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodGet {
		matchId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
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
			log.Println("WARN: error decoding bookmark to add", err)
			return
		}

		if !insData.Archived {
			db.Add(database, insData)
		} else {
			insData.Archived, insData.SnapshotURL = archive.SendSnapshot(insData.URL)
			if insData.Archived {
				db.Add(database, insData)
			} else {
				log.Println("WARN: Snapshot failed:", insData.URL)
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

func updateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodPost {
		matchId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			return
		}

		var newData setup.Bookmark
		err = json.NewDecoder(r.Body).Decode(&newData)
		if err != nil {
			log.Println("WARN: error decoding bookmark to update", err)
			return
		}
		newData.ID = matchId

		if !newData.Archived {
			newData.SnapshotURL = ""
			db.Update(database, newData, true)
		} else if newData.Archived && newData.SnapshotURL != "" {
			db.Update(database, newData, true)
		} else {
			newData.Archived, newData.SnapshotURL = archive.SendSnapshot(newData.URL)
			if newData.Archived {
				db.Update(database, newData, true)
			} else {
				log.Println("WARN: Snapshot failed", newData.URL)
				db.Update(database, newData, true)
			}
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
	log.Printf("WARN: status 500 at '%s%s'\n", r.Host, r.URL)
}

func rowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		matchId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			return
		}

		oldData, err := db.ViewSingleRow(database, matchId)
		if err != nil {
			log.Printf("WARN: could not get record for bookmark id: %d: %v\n", matchId, err)
			return
		}

		writeData, err := json.Marshal(&oldData)
		if err != nil {
			log.Println("WARN: could not marshal bookmark data:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(writeData))
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
			log.Println("WARN: failed to decode bookmark data:", err)
			return
		}

		var searchUrl string = getData.URL
		getData, err = db.SearchByUrl(database, searchUrl)
		if err != nil {
			log.Println("WARN: failed searching by URL:", err)
			return
		}

		if getData.ID == 0 {
			http.NotFound(w, r)
			return
		}

		writeData, err := json.Marshal(&getData)
		if err != nil {
			log.Println("WARN: could not marshal bookmark data:", err)
		}
		w.Write([]byte(writeData))
		db.AddToRecents(database, getData.ID)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func refetchThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodGet {
		matchId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			return
		}

		if err := db.RefetchThumbnail(database, matchId, nil); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Thumbnail updated"))
	} else if r.Method == http.MethodPost {
		matchId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			log.Println("WARN: invalid bookmark id:", err)
			return
		}

		if err := r.ParseMultipartForm(constants.THUMBNAIL_FILE_SIZE); err != nil {
			http.Error(w, "error parsing data", http.StatusBadRequest)
			log.Println("WARN: could not parse thumbnail data:", err)
			return
		}

		thumbnailFile, _, err := r.FormFile("thumbnail")
		if err != nil {
			http.Error(w, "error getting thumbnail field", http.StatusBadRequest)
			log.Println("ERROR: getting thumbnail field:", err)
			return
		}
		defer thumbnailFile.Close()

		thumbnailFileBytes, err := io.ReadAll(thumbnailFile)
		if err != nil {
			http.Error(w, "error converting thumbnail file to bytes.", http.StatusInternalServerError)
			log.Println("ERROR: converting thumbnail file to bytes:", err)
			return
		}

		if err := db.RefetchThumbnail(database, matchId, thumbnailFileBytes); err != nil {
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

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodGet {
		listCategories, err := db.GetAllCategories(database)
		if err != nil {
			log.Println("ERROR: getting categories:", err)
			return
		}

		for _, category := range listCategories {
			fmt.Fprintf(w, "<option value=\"%s\"></option>", category)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func pagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprint(w, db.TotalPageCount(database))
	} else if r.Method == http.MethodPost {
		fmt.Fprint(w, pageCountForSearch)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func Start(data *sql.DB) {
	database = data

	mux := http.NewServeMux()

	fsopen := fs.FS(Web)
	webStatic, err := fs.Sub(fsopen, "web/static")
	if err != nil {
		log.Fatalln("ERROR: could not open embedded 'web' directory:", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(webStatic))))

	tmplFuncMap["getHostname"] = getHostname
	tmplFuncMap["keywordSplit"] = keywordSplit
	tmplFuncMap["grabThumbnail"] = grabThumbnail
	tmplFuncMap["PageTitle"] = pageTitle
	tmplFuncMap["encapsulateURL"] = encapsulateURL

	mux.HandleFunc("/{$}", rootHandler)
	mux.HandleFunc("/import/", importHandler)
	mux.HandleFunc("/api/import-bookmark/", importBookmarkHandler)
	mux.HandleFunc("/api/add/", addHandler)
	mux.HandleFunc("/api/row/{id}", rowHandler)
	mux.HandleFunc("/api/update/{id}", updateHandler)
	mux.HandleFunc("/api/delete/{id}", deleteHandler)
	mux.HandleFunc("/api/categories/", categoriesHandler)
	mux.HandleFunc("/api/check-url/", checkUrlHandler)
	mux.HandleFunc("/api/refetch-thumbnail/{id}", refetchThumbnailHandler)
	mux.HandleFunc("/api/pages/", pagesHandler)

	go secondaryServer()

	log.Println("INFO: Web-server starting at:", constants.WEBUI_ADDR)
	if err := http.ListenAndServe(constants.WEBUI_ADDR, mux); err != nil {
		log.Fatalln("ERROR: could not start web UI server:", err)
	}
}

func secondaryServer() {
	secondaryMux := http.NewServeMux()
	secondaryMux.Handle("GET /thumbnail/", http.StripPrefix("/thumbnail/", http.FileServer(http.Dir(constants.THUMBNAILS_PATH))))
	// secondaryMux.Handle("GET /archive/", http.StripPrefix("/archive/", http.FileServer(http.Dir(constants.ARCHIVES_PATH))))

	secondaryMux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Alive")
	})

	if err := http.ListenAndServe(constants.SECONDARY_PORT, secondaryMux); err != nil {
		log.Fatalln("ERROR: could not start secondary server:", err)
	}
}

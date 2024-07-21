package server

import (
	"dalennod/internal/archive"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const PORT string = ":41415"

var (
	tmplFuncMap template.FuncMap = make(template.FuncMap)
	database    *sql.DB
	tmpl        *template.Template
	Web         embed.FS
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

	tmplFuncMap["keywordSplit"] = keywordSplit

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/delete/", deleteHandler)
	mux.HandleFunc("/add/", addHandler)
	mux.HandleFunc("/row/", rowHandler)
	mux.HandleFunc("/groups/", groupsHandler)
	mux.HandleFunc("/update/", updateHandler)
	mux.HandleFunc("/search/", searchHandler)
	mux.HandleFunc("/searchKeyword/", searchKeywordHandler)
	mux.HandleFunc("/searchGroup/", searchGroupHandler)
	mux.HandleFunc("/checkUrl/", checkUrlHandler)

	logger.Info.Printf("Web-server starting on http://localhost%s/\n", PORT)
	fmt.Printf("Web-server starting on http://localhost%s/\n", PORT)
	// TODO: use own Server to gracefully shutdown
	err = http.ListenAndServe(PORT, mux)
	if err != nil {
		fmt.Printf("Stopping (error: %v)\n", err)
		logger.Error.Printf("Stopping (error: %v)\n", err)
	}
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
	logger.Warn.Printf("status 500 at '%s%s'\n", r.Host, r.URL)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
	logger.Warn.Printf("status 404 at '%s%s'\n", r.Host, r.URL)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl = template.Must(template.New("index").Funcs(tmplFuncMap).ParseFS(Web, "web/index.html"))

		var bookmarks []setup.Bookmark = db.ViewAll(database, true)
		if err := execTmplBmInterface(w, "index", bookmarks); err != nil {
			logger.Warn.Println(err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == "GET" {
		var (
			deleteID *regexp.Regexp = regexp.MustCompile("^/delete/([0-9]+)$")
			match    []string       = deleteID.FindStringSubmatch(r.URL.Path)
		)
		if len(match) < 2 {
			internalServerErrorHandler(w, r)
			return
		}
		matchInt, err := strconv.Atoi(match[1])
		if err != nil {
			logger.Error.Println(err)
		}

		db.Remove(database, matchInt)

		rootHandler(w, r)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == "POST" {
		var insData setup.Bookmark
		var err error = json.NewDecoder(r.Body).Decode(&insData)
		if err != nil {
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
	} else if r.Method == "GET" {
		w.Write([]byte("Alive"))
	} else {
		internalServerErrorHandler(w, r)
	}
}

func rowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			oldData setup.Bookmark
			rowID   *regexp.Regexp = regexp.MustCompile("^/row/([0-9]+)$")
			match   []string       = rowID.FindStringSubmatch(r.URL.Path)
		)

		if len(match) < 2 {
			internalServerErrorHandler(w, r)
			return
		}
		matchInt, err := strconv.Atoi(match[1])
		if err != nil {
			logger.Error.Println(err)
			return
		}

		oldData, err = db.ViewSingleRow(database, matchInt, true)
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
	if r.Method == "POST" {
		var (
			newData  setup.Bookmark
			updateID *regexp.Regexp = regexp.MustCompile("^/update/([0-9]+)$")

			match []string = updateID.FindStringSubmatch(r.URL.Path)
		)
		if len(match) < 2 {
			internalServerErrorHandler(w, r)
			return
		}
		matchInt, err := strconv.Atoi(match[1])
		if err != nil {
			logger.Error.Println(err)
		}

		err = json.NewDecoder(r.Body).Decode(&newData)
		if err != nil {
			logger.Error.Println(err)
		}
		newData.ID = matchInt

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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == "POST" {
		r.ParseForm()
		var searchTerm string = r.FormValue("searchTerm")

		if searchTerm == "" {
			if err := execTmplBmInterface(w, "bm_list", db.ViewAll(database, true)); err != nil {
				logger.Error.Println(err)
			}
			return
		}

		var bookmarks []setup.Bookmark = db.ViewAllWhere(database, searchTerm)

		// TODO: in future look into htmx pop-up on 404 status
		// if len(bookmarks) == 0 {
		// 	notFoundHandler(w, r)
		// 	return
		// }

		if err := execTmplBmInterface(w, "bm_list", bookmarks); err != nil {
			logger.Error.Println(err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func searchKeywordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == "POST" {
		r.ParseForm()
		var searchTerm string = r.FormValue("searchTerm")

		var bookmarks []setup.Bookmark = db.ViewAllWhereKeyword(database, searchTerm)
		if err := execTmplBmInterface(w, "bm_list", bookmarks); err != nil {
			logger.Error.Println(err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func searchGroupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == "POST" {
		r.ParseForm()
		var searchTerm string = r.FormValue("searchTerm")

		var bookmarks []setup.Bookmark = db.ViewAllWhereGroup(database, searchTerm)
		if err := execTmplBmInterface(w, "bm_list", bookmarks); err != nil {
			logger.Error.Println(err)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}

func checkUrlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == "POST" {
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
			notFoundHandler(w, r)
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

func keywordSplit(keywords string, delimiter string) []string {
	return strings.Split(keywords, delimiter)
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == "GET" {
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

func execTmplBmInterface(wr io.Writer, name string, data any) error {
	var allBookmarks map[string]interface{} = make(map[string]interface{})
	allBookmarks["Bookmarks"] = data
	if err := tmpl.ExecuteTemplate(wr, name, allBookmarks); err != nil {
		return err
	}
	return nil
}

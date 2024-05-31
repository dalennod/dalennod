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
	"io/fs"
	"net/http"
	"regexp"
	"strconv"
)

const PORT string = ":41415"

var (
	data *sql.DB
	tmpl *template.Template

	IndexHtml embed.FS
	Webui     embed.FS
)

func Start(database *sql.DB) {
	data = database
	var mux *http.ServeMux = http.NewServeMux()

	fsopen := fs.FS(Webui)
	webuiStatic, _ := fs.Sub(fsopen, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(webuiStatic))))

	mux.HandleFunc("/", root)
	mux.HandleFunc("/delete/", deleteHandler)
	mux.HandleFunc("/add/", addHandler)
	mux.HandleFunc("/getRow/", getRowHandler)
	mux.HandleFunc("/update/", updateHandler)
	mux.HandleFunc("/static/search.html", searchHandler)
	mux.HandleFunc("/checkUrl/", checkUrlHandler)

	logger.Info.Printf("Web-server starting on http://localhost%s\n", PORT)
	fmt.Printf("Web-server starting on http://localhost%s\n", PORT)
	err := http.ListenAndServe(PORT, mux)
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

func root(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl = template.Must(template.ParseFS(IndexHtml, "index.html"))
		var bookmarks []setup.Bookmark = db.ViewAll(data, true)
		tmpl.Execute(w, bookmarks)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == "GET" {
		var (
			deleteID          = regexp.MustCompile("^/delete/([0-9]+)$")
			match    []string = deleteID.FindStringSubmatch(r.URL.Path)
		)
		if len(match) < 2 {
			internalServerErrorHandler(w, r)
			return
		}
		matchInt, err := strconv.Atoi(match[1])
		if err != nil {
			logger.Error.Println(err)
		}

		db.Remove(data, matchInt)
		w.WriteHeader(http.StatusOK)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == "POST" {
		var (
			insData       setup.Bookmark
			archiveResult bool   = false
			snapshotURL   string = ""
		)
		err := json.NewDecoder(r.Body).Decode(&insData)
		if err != nil {
			logger.Error.Println(err)
		}

		if !insData.Archived {
			db.Add(data, insData.URL, insData.Title, insData.Note, insData.Keywords, insData.BGroup, insData.Archived, snapshotURL, "s")
		} else {
			archiveResult, snapshotURL = archive.SendSnapshot(insData.URL)
			if archiveResult {
				db.Add(data, insData.URL, insData.Title, insData.Note, insData.Keywords, insData.BGroup, insData.Archived, snapshotURL, "s")
			} else {
				logger.Warn.Println("Snapshot failed.", snapshotURL)
				db.Add(data, insData.URL, insData.Title, insData.Note, insData.Keywords, insData.BGroup, false, snapshotURL, "s")
			}
		}
		w.WriteHeader(http.StatusCreated)
	} else if r.Method == "GET" {
		w.Write([]byte("Alive."))
	} else {
		internalServerErrorHandler(w, r)
	}
}

func getRowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var (
			oldData  setup.Bookmark
			getRowID          = regexp.MustCompile("^/getRow/([0-9]+)$")
			match    []string = getRowID.FindStringSubmatch(r.URL.Path)
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
		oldData.ID = matchInt

		oldData.URL, oldData.Title, oldData.Note, oldData.Keywords, oldData.BGroup, oldData.Archived = db.ViewSingleRow(data, matchInt, true)

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
			newData       setup.Bookmark
			archiveResult bool     = false
			snapshotURL   string   = ""
			updateID               = regexp.MustCompile("^/update/([0-9]+)$")
			match         []string = updateID.FindStringSubmatch(r.URL.Path)
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
			db.Update(data, newData.URL, newData.Title, newData.Note, newData.Keywords, newData.BGroup, newData.ID, false, true, snapshotURL)
		} else {
			archiveResult, snapshotURL = archive.SendSnapshot(newData.URL)
			if archiveResult {
				db.Update(data, newData.URL, newData.Title, newData.Note, newData.Keywords, newData.BGroup, newData.ID, true, true, snapshotURL)
			} else {
				logger.Warn.Println("Snapshot failed.")
				db.Update(data, newData.URL, newData.Title, newData.Note, newData.Keywords, newData.BGroup, newData.ID, false, true, snapshotURL)
			}
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		var bookmarks []setup.Bookmark
		for _, vs := range r.Form {
			for _, v := range vs {
				bookmarks = db.ViewAllWhere(data, v)
			}
		}

		if len(bookmarks) == 0 {
			notFoundHandler(w, r)
			return
		}

		tmpl = template.Must(template.ParseFS(Webui, "static/search.html"))
		tmpl.Execute(w, bookmarks)
	} else {
		internalServerErrorHandler(w, r)
	}
}

func checkUrlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "POST" {
		var getData setup.Bookmark
		var err error = json.NewDecoder(r.Body).Decode(&getData)
		if err != nil {
			logger.Error.Println(err)
		}

		var searchUrl string = getData.URL
		getData = db.SearchByUrl(data, searchUrl)

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

package server

import (
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"net/http"
)

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

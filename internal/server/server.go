package server

import (
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"database/sql"
	"html/template"
	"net/http"
)

const PORT string = "41415"

var (
	data *sql.DB
	tmpl *template.Template
)

func Start(database *sql.DB) {
	data = database
	tmpl = template.Must(template.ParseFiles("index.html"))

	var fileServer http.Handler = http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.HandleFunc("/", root)
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		logger.Error.Println(err)
	}
	logger.Info.Println("Started on http://localhost:" + PORT)
}

func root(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte(fmt.Sprintf("<h1>%s</h1>", db.ViewAll(data, "s"))))
	var bookmarks []db.Bookmark = db.ViewAll(data, "s")
	tmpl.Execute(w, bookmarks)
}

package server

import (
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"database/sql"
	"fmt"
	"net/http"
)

const PORT string = "41415"

var data *sql.DB

func Start(database *sql.DB) {
	data = database

	http.HandleFunc("/", root)
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		logger.Error.Println(err)
	}
	logger.Info.Println("Started on http://localhost:" + PORT)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("<h1>%s</h1>", db.ViewAll(data, "s"))))
}

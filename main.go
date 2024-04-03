package main

import (
	"dalennod/internal/server"
	"dalennod/internal/setup"
	"dalennod/internal/user_input"
	"embed"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed index.html
	indexHtml embed.FS
	//go:embed static
	webui embed.FS
)

func init() {
	server.IndexHtml = indexHtml
	server.Webui = webui
}

func main() {
	user_input.UserInput(setup.CreateDB(setup.GetOS()))
}

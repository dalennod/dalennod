package main

import (
	"dalennod/internal/server"
	"dalennod/internal/setup"
	"dalennod/internal/user_input"
	"embed"

	_ "github.com/mattn/go-sqlite3" // CGO driver
	// _ "modernc.org/sqlite" // CGO-free driver
)

var (
	//go:embed web
	web embed.FS
)

func init() {
	server.Web = web
}

func main() {
	user_input.UserInput(setup.CreateDB(setup.InitLocalDirs()))
}

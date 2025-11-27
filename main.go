package main

import (
	"embed"
	"flag"
	"log"
	"os"

	"dalennod/internal/server"
	"dalennod/internal/setup"
	"dalennod/internal/user_input"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed web
var web embed.FS

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	server.Web = web
	setup.ParseFlags()
}

func main() {
	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(0)
	}
	user_input.UserInput(setup.CreateDB(setup.InitLocalDirs()))
}

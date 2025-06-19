package main

import (
    "embed"
    "flag"
    "os"

    "dalennod/internal/server"
    "dalennod/internal/setup"
    "dalennod/internal/user_input"

    _ "github.com/mattn/go-sqlite3" // CGO driver
    // _ "modernc.org/sqlite" // CGO-free driver
)

//go:embed web
var web embed.FS

func init() {
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

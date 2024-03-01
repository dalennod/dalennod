package main

import (
	"dalennod/internal/setup"
	"dalennod/internal/user_input"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	user_input.UserInput(setup.CreateDB(setup.GetOS()))
}

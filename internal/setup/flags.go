package setup

import "flag"

type FlagValues struct {
	RemoveID    string
	UpdateID    string
	ViewID      string
	ViewAll     bool
	AddEntry    bool
	StartServer bool
	Backup      bool
	JSONOut     bool
}

var flagValues FlagValues

func cliFlags() {
	flag.StringVar(&flagValues.RemoveID, "r", "", "Delete specific bookmark using its ID.")
	flag.StringVar(&flagValues.RemoveID, "remove", "", "Delete specific bookmark using its ID.")

	flag.StringVar(&flagValues.UpdateID, "u", "", "Update specific bookmark using its ID.")
	flag.StringVar(&flagValues.UpdateID, "update", "", "Update specific bookmark using its ID.")

	flag.StringVar(&flagValues.ViewID, "v", "", "View specific bookmark using its ID.")
	flag.StringVar(&flagValues.ViewID, "view", "", "View specific bookmark using its ID.")

	flag.BoolVar(&flagValues.ViewAll, "va", false, "View all bookmarks.")
	flag.BoolVar(&flagValues.ViewAll, "view-all", false, "View all bookmarks.")

	flag.BoolVar(&flagValues.AddEntry, "a", false, "Add a bookmark entry to the database.")
	flag.BoolVar(&flagValues.AddEntry, "add", false, "Add a bookmark entry to the database.")

	flag.BoolVar(&flagValues.StartServer, "s", false, "Start webserver locally for UI.")
	flag.BoolVar(&flagValues.StartServer, "serve", false, "Start webserver locally for UI.")

	flag.BoolVar(&flagValues.Backup, "b", false, "Start backup process.")
	flag.BoolVar(&flagValues.Backup, "backup", false, "Start backup process.")
	flag.BoolVar(&flagValues.JSONOut, "json", false, "Print entire DB in JSON. Use alongside --backup flag.")
}

func ParseFlags() FlagValues {
	cliFlags()
	flag.Parse()
	return flagValues
}

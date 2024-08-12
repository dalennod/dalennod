package setup

import (
	"flag"
	"fmt"
	"io"
)

type FlagValues struct {
	RemoveID    string
	UpdateID    string
	ViewID      string
	ViewAll     bool
	AddEntry    bool
	StartServer bool
	Backup      bool
	JSONOut     bool
	Import      string
	Firefox     bool
	Dalennod    bool
	Where       bool
}

var flagValues FlagValues

func cliFlags() {
	flag.Usage = func() {
		var w io.Writer = flag.CommandLine.Output()
		fmt.Fprintln(w, "Usage of dalennod: dalennod [OPTION] ...")
		fmt.Fprintln(w, "\nOptions:")
		fmt.Fprintln(w, "  -s, --serve\t\tStart webserver locally for Web UI & Extension")
		fmt.Fprintln(w, "  -a, --add\t\tAdd a bookmark entry to the database")
		fmt.Fprintln(w, "  -r, --remove [id]\tRemove specific bookmark using its ID")
		fmt.Fprintln(w, "  -u, --update [id]\tUpdate specific bookmark using its ID")
		fmt.Fprintln(w, "  -v, --view [id]\tView specific bookmark using its ID")
		fmt.Fprintln(w, "  -va, --view-all\tView all bookmarks")
		fmt.Fprintln(w, "  -i, --import [file]\tImport bookmarks from a browser")
		fmt.Fprintln(w, "  -ff, --firefox\tImport bookmarks from Firefox \n\t\t\t  Must use alongside -i, --import option")
		fmt.Fprintln(w, "  -41, --dalennod\tImport bookmarks from exported Dalennod JSON \n\t\t\t  Must use alongside -i, --import option")
		fmt.Fprintln(w, "  -b, --backup\t\tStart backup process")
		fmt.Fprintln(w, "  --json\t\tPrint entire DB in JSON \n\t\t\t  Use alongside -b, --backup flag")
		fmt.Fprintln(w, "  --where\t\tPrint config and logs directory path")
		fmt.Fprintln(w, "  -h, --help\t\tShows this message")
	}

	flag.BoolVar(&flagValues.StartServer, "s", false, "Start webserver locally for Web UI & Extension")
	flag.BoolVar(&flagValues.StartServer, "serve", false, "Start webserver locally for Web UI & Extension")

	flag.StringVar(&flagValues.RemoveID, "r", "", "Remove specific bookmark using its ID")
	flag.StringVar(&flagValues.RemoveID, "remove", "", "Remove specific bookmark using its ID")

	flag.StringVar(&flagValues.UpdateID, "u", "", "Update specific bookmark using its ID")
	flag.StringVar(&flagValues.UpdateID, "update", "", "Update specific bookmark using its ID")

	flag.StringVar(&flagValues.ViewID, "v", "", "View specific bookmark using its ID")
	flag.StringVar(&flagValues.ViewID, "view", "", "View specific bookmark using its ID")

	flag.BoolVar(&flagValues.ViewAll, "va", false, "View all bookmarks")
	flag.BoolVar(&flagValues.ViewAll, "view-all", false, "View all bookmarks")

	flag.BoolVar(&flagValues.AddEntry, "a", false, "Add a bookmark entry to the database")
	flag.BoolVar(&flagValues.AddEntry, "add", false, "Add a bookmark entry to the database")

	flag.BoolVar(&flagValues.Backup, "b", false, "Start backup process")
	flag.BoolVar(&flagValues.Backup, "backup", false, "Start backup process")
	flag.BoolVar(&flagValues.JSONOut, "json", false, "Print entire DB in JSON. Use alongside --backup flag")

	flag.StringVar(&flagValues.Import, "i", "", "Import bookmarks from a browser")
	flag.StringVar(&flagValues.Import, "import", "", "Import bookmarks from a browser")
	flag.BoolVar(&flagValues.Firefox, "ff", false, "Import bookmarks from Firefox. Use alongside -i flag")
	flag.BoolVar(&flagValues.Firefox, "firefox", false, "Import bookmarks from Firefox. Use alongside -i flag")
	flag.BoolVar(&flagValues.Dalennod, "41", false, "Import bookmarks exported Dalennod JSON. Use alongside -i flag")
	flag.BoolVar(&flagValues.Dalennod, "dalennod", false, "Import bookmarks exported Dalennod JSON. Use alongside -i flag")

	flag.BoolVar(&flagValues.Where, "where", false, "Print config and logs directory path")
}

func ParseFlags() FlagValues {
	cliFlags()
	flag.Parse()
	return flagValues
}

package setup

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	Import      bool
	Firefox     string
	Chromium    string
	Dalennod    string
	Where       bool
	Profile     bool
	Switch      string
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
		fmt.Fprintln(w, "  -V, --view-all\tView all bookmarks")
		fmt.Fprintln(w, "  -i, --import [file]\tImport bookmarks from a browser")
		fmt.Fprintln(w, "  --firefox\t\tImport bookmarks from Firefox \n\t\t\t  Must use alongside -i, --import option")
		fmt.Fprintln(w, "  --chromium\t\tImport bookmarks from Chromium \n\t\t\t  Must use alongside -i, --import option")
		fmt.Fprintln(w, "  --dalennod\t\tImport bookmarks from exported Dalennod JSON \n\t\t\t  Must use alongside -i, --import option")
		fmt.Fprintln(w, "  -b, --backup\t\tStart backup process")
		fmt.Fprintln(w, "  --json\t\tPrint entire DB in JSON \n\t\t\t  Use alongside -b, --backup flag")
		fmt.Fprintln(w, "  --where\t\tPrint config and logs directory path")
		fmt.Fprintln(w, "  --profile\t\tShow profile names found in local directory")
		fmt.Fprintln(w, "  --switch\t\tSwitch profiles \n\t\t\t Must use alongside --profile flag")
		fmt.Fprintln(w, "  -h, --help\t\tShows help message")
	}

	flag.BoolVar(&flagValues.StartServer, "s", false, "Start webserver locally for Web UI & Extension")
	flag.BoolVar(&flagValues.StartServer, "serve", false, "Start webserver locally for Web UI & Extension")

	flag.StringVar(&flagValues.RemoveID, "r", "", "Remove specific bookmark using its ID")
	flag.StringVar(&flagValues.RemoveID, "remove", "", "Remove specific bookmark using its ID")

	flag.StringVar(&flagValues.UpdateID, "u", "", "Update specific bookmark using its ID")
	flag.StringVar(&flagValues.UpdateID, "update", "", "Update specific bookmark using its ID")

	flag.StringVar(&flagValues.ViewID, "v", "", "View specific bookmark using its ID")
	flag.StringVar(&flagValues.ViewID, "view", "", "View specific bookmark using its ID")

	flag.BoolVar(&flagValues.ViewAll, "V", false, "View all bookmarks")
	flag.BoolVar(&flagValues.ViewAll, "view-all", false, "View all bookmarks")

	flag.BoolVar(&flagValues.AddEntry, "a", false, "Add a bookmark entry to the database")
	flag.BoolVar(&flagValues.AddEntry, "add", false, "Add a bookmark entry to the database")

	flag.BoolVar(&flagValues.Backup, "b", false, "Start backup process")
	flag.BoolVar(&flagValues.Backup, "backup", false, "Start backup process")
	flag.BoolVar(&flagValues.JSONOut, "json", false, "Print entire DB in JSON. Use alongside --backup flag")

	flag.BoolVar(&flagValues.Import, "i", false, "Import bookmarks from a browser")
	flag.BoolVar(&flagValues.Import, "import", false, "Import bookmarks from a browser")
	flag.StringVar(&flagValues.Firefox, "firefox", "", "Import bookmarks from Firefox. Use alongside -i flag")
	flag.StringVar(&flagValues.Chromium, "chromium", "", "Import bookmarks from Chromium. Use alongside -i flag")
	flag.StringVar(&flagValues.Dalennod, "dalennod", "", "Import bookmarks exported Dalennod JSON. Use alongside -i flag")

	flag.BoolVar(&flagValues.Where, "where", false, "Print config and logs directory path")

	flag.BoolVar(&flagValues.Profile, "profile", false, "Show profile names found in local directory")
	flag.StringVar(&flagValues.Switch, "switch", "", "Switch profiles. Must use alongside --profile flag")
}

func ParseFlags() FlagValues {
	cliFlags()
	flag.Parse()
	return flagValues
}

func setCompletion() {
	shell := os.Getenv("SHELL")
	switch {
	case strings.Contains(shell, "fish"):
		fishCompletion()
	case strings.Contains(shell, "bash"):
		bashCompletion()
	case strings.Contains(shell, "zsh"):
		zshCompletion()
	}
}

func fishCompletion() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Println("error finding home directory. ERROR:", err)
		return
	}

	fishLocalPath := filepath.Join(homePath, ".local", "share", "fish", "generated_completions")
	fishLocalStat, err := os.Stat(fishLocalPath)
	if err != nil {
		log.Println("error getting fish shell local directory info. ERROR", err)
		return
	}

	if !fishLocalStat.IsDir() {
		os.MkdirAll(fishLocalPath, 0755)
	}

	fishCompletionPath := filepath.Join(fishLocalPath, "dalennod.fish")
	if _, err := os.Stat(fishCompletionPath); os.IsExist(err) {
		return
	}

	fishCompletionFile, err := os.Create(fishCompletionPath)
	if err != nil {
		log.Println("error creating fish completion file. ERROR:", err)
	}
	defer fishCompletionFile.Close()

	_, err = fishCompletionFile.Write([]byte(`# Autogenerated from dalennod program
# $HOME/.local/share/fish/generated_completions/dalennod.fish

complete -c dalennod -s a -l add -d "Add a bookmark entry to the database"
complete -c dalennod -s V -l view-all -d "View all bookmarks"
complete -c dalennod -s v -l view -r -d "View specific bookmark using its ID"
complete -c dalennod -s u -l update -r -d "Update specific bookmark using its ID"
complete -c dalennod -s r -l remove -r -d "Remove specific bookmark using its ID"
complete -c dalennod -s b -l backup -d "Start backup process"
complete -c dalennod -s s -l serve -d "Start webserver locally for Web UI & Extension"
complete -c dalennod -s h -l help -d "Shows help message"
complete -c dalennod -l json -d "Print entire DB in JSON. Use alongside -b, --backup flag"
complete -c dalennod -l where -d "Print config and logs directory path"
complete -c dalennod -l profile -d "Show profile names found in local directory"
complete -c dalennod -l switch -d "Switch profiles. Must use alongside --profile flag"

# import options
complete -c dalennod -s i -l import -d "Import bookmarks from a browser"
complete -c dalennod -l firefox -d "Import bookmarks from Firefox. Must use alongside -i, --import option"
complete -c dalennod -l chromium -d "Import bookmarks from Chromium. Must use alongside -i, --import option"
complete -c dalennod -l dalennod -d "Import bookmarks from exported Dalennod JSON. Must use alongside -i, --import option"
`))
	if err != nil {
		log.Println("error writing to fish completion file. ERROR:", err)
	}
}

func bashCompletion() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Println("error finding home directory. ERROR:", err)
	}

	bashrcPath := filepath.Join(homePath, ".bashrc")
	bashrcFile, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("error opening bashrc. ERROR:", bashrcPath)
	}
	defer bashrcFile.Close()

	_, err = bashrcFile.Write([]byte(`# Autogenerated from dalennod program

_dalennod() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="-s --serve -a --add -r --remove -u --update -v --view -V --view-all -i --import --firefox --chromium --dalennod -b --backup --json --where --profile --switch"

    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi
}

complete -F _dalennod dalennod

# End of autogenerated lines from dalennod program
`))
	if err != nil {
		log.Println("error writing to bashrc. ERROR:", err)
	}
}

func zshCompletion() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Println("error finding home directory. ERROR:", err)
	}

	configDir, err := ConfigDir()
	if err != nil {
		log.Println("error getting config directory. ERROR:", err)
	}

	compdefPath := filepath.Join(configDir, "zsh-completion")
	err = os.MkdirAll(compdefPath, 0755)
	if err != nil {
		log.Println("error creating zsh-completion directory inside config dir. ERROR:", err)
	}

	compdefDalennodPath := filepath.Join(compdefPath, "_dalennod")
	compdefDalennodFile, err := os.Create(compdefDalennodPath)
	if err != nil {
		log.Println("error creating compdef _dalennod file. ERROR:", err)
	}
	defer compdefDalennodFile.Close()

	_, err = compdefDalennodFile.Write([]byte(`#compdef dalennod
# Autogenerated from dalennod program

_arguments -s \
    '(-s,--serve)'{-s,--serve}'[Start webserver locally for Web UI & Extension]' \
    '(-a,--add)'{-a,--add}'[Add a bookmark entry to the database]' \
    '(-r,--remove)'{-r,--remove}'[Remove specific bookmark using its ID]' \
    '(-u,--update)'{-u,--update}'[Update specific bookmark using its ID]' \
    '(-v,--view)'{-v,--view}'[View specific bookmark using its ID]' \
    '(-V,--view-all)'{-V,--view-all}'[View all bookmarks]' \
    '(-i,--import)'{-i,--import}'[Import bookmarks from a browser]' \
    '(-h,--help)'{-h,--help}'[Shows help message]' \
    '--firefox[Import bookmarks from Firefox. Must use alongside -i, --import option]' \
    '--chromium[Import bookmarks from Chromium. Must use alongside -i, --import option]' \
    '--dalennod[Import bookmarks exported Dalennod JSON. Must use alongside -i, --import option]' \
    '(-b,--backup)'{-b,--backup}'[Start backup process]' \
    '--json[Print entire DB in JSON. Use alongside -b, --backup flag]' \
    '--where[Print config and logs directory path]'
    '--profile[Show profile names found in local directory]' \
    '--switch[Switch profiles. Must use alongside --profile flag]'
`))
	if err != nil {
		log.Println("error writing to compdef _dalennod file. ERROR:", err)
	}

	zshrcPath := filepath.Join(homePath, ".zshrc")
	zshrcFile, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("error opening bashrc. ERROR:", zshrcPath)
	}
	defer zshrcFile.Close()

	zshrcFileWrite := "# Autogenerated from dalennod program\n"
	zshrcFileWrite += "export FPATH=" + compdefPath + ":$FPATH\n"
	zshrcFileWrite += "autoload -Uz compinit && compinit\n"
	zshrcFileWrite += "# End of autogenerated lines from dalennod program\n"

	_, err = zshrcFile.Write([]byte(zshrcFileWrite))
	if err != nil {
		log.Println("error writing into .zshrc. ERROR:", err)
	}
}

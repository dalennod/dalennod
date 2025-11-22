package setup

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"dalennod/internal/constants"
)

type FlagValues struct {
	RemoveID       string
	UpdateID       string
	ViewID         string
	ViewAll        bool
	AddEntry       bool
	StartServer    bool
	Backup         bool
	JSONOut        bool
	Crypt          bool
	RedoCompletion bool
	Import         bool
	Firefox        string
	Chromium       string
	Dalennod       string
	Where          bool
	Profile        bool
	Switch         string
	FixDB          bool
}

var FlagVals FlagValues

func ParseFlags() FlagValues {
	cliFlags()
	flag.Parse()
	return FlagVals
}

func cliFlags() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of dalennod (%s): dalennod [OPTION] ...\n\n", constants.VERSION)
		fmt.Fprintln(w, "Options:")
		fmt.Fprintln(w, "  -s, --serve         Start webserver locally for Web UI & Extension")
		fmt.Fprintln(w, "  -a, --add           Add a bookmark entry to the database")
		fmt.Fprintln(w, "  -r, --remove [id]   Remove specific bookmark using its ID")
		fmt.Fprintln(w, "  -u, --update [id]   Update specific bookmark using its ID")
		fmt.Fprintln(w, "  -v, --view   [id]   View specific bookmark using its ID")
		fmt.Fprintln(w, "  -V, --view-all      View all bookmarks")
		fmt.Fprintln(w, "  -i, --import        Import bookmarks from a browser")
		fmt.Fprintln(w, "  --firefox  [file]   Import bookmarks from Firefox")
		fmt.Fprintln(w, "                        Must use alongside -i, --import option")
		fmt.Fprintln(w, "  --chromium [file]   Import bookmarks from Chromium")
		fmt.Fprintln(w, "                        Must use alongside -i, --import option")
		fmt.Fprintln(w, "  --dalennod [file]   Import bookmarks from exported Dalennod JSON")
		fmt.Fprintln(w, "                        Must use alongside -i, --import option")
		fmt.Fprintln(w, "  -b, --backup        Start backup process")
		fmt.Fprintln(w, "  --json              Dump entire DB in JSON")
		fmt.Fprintln(w, "                        Use alongside -b, --backup flag")
		fmt.Fprintln(w, "  --crypt             Encrypt/decrypt the JSON backup")
		fmt.Fprintln(w, "                        Use alongside --json flag to encrypt")
		fmt.Fprintln(w, "                        Use alongside --import --dalennod to decrypt")
		fmt.Fprintln(w, "  --where             Print config and logs directory location")
		fmt.Fprintln(w, "  --redo-completion   Set command line completion again")
		fmt.Fprintln(w, "  --profile           Show profile names found in local directory")
		fmt.Fprintln(w, "  --switch [profile]  Switch profiles")
		fmt.Fprintln(w, "                        Must use alongside --profile flag")
		fmt.Fprintln(w, "  -h, --help          Shows this help message")
	}

	flag.BoolVar(&FlagVals.StartServer, "s", false, "Start webserver locally for Web UI & Extension")
	flag.BoolVar(&FlagVals.StartServer, "serve", false, "Start webserver locally for Web UI & Extension")

	flag.StringVar(&FlagVals.RemoveID, "r", "", "Remove specific bookmark using its ID")
	flag.StringVar(&FlagVals.RemoveID, "remove", "", "Remove specific bookmark using its ID")

	flag.StringVar(&FlagVals.UpdateID, "u", "", "Update specific bookmark using its ID")
	flag.StringVar(&FlagVals.UpdateID, "update", "", "Update specific bookmark using its ID")

	flag.StringVar(&FlagVals.ViewID, "v", "", "View specific bookmark using its ID")
	flag.StringVar(&FlagVals.ViewID, "view", "", "View specific bookmark using its ID")

	flag.BoolVar(&FlagVals.ViewAll, "V", false, "View all bookmarks")
	flag.BoolVar(&FlagVals.ViewAll, "view-all", false, "View all bookmarks")

	flag.BoolVar(&FlagVals.AddEntry, "a", false, "Add a bookmark entry to the database")
	flag.BoolVar(&FlagVals.AddEntry, "add", false, "Add a bookmark entry to the database")

	flag.BoolVar(&FlagVals.Backup, "b", false, "Start backup process")
	flag.BoolVar(&FlagVals.Backup, "backup", false, "Start backup process")
	flag.BoolVar(&FlagVals.JSONOut, "json", false, "Dump entire DB in JSON. Use alongside --backup flag")
	flag.BoolVar(&FlagVals.Crypt, "crypt", false, "Encrypt/decrypt the JSON backup. Use alongside --json flag to encrypt or alongside --import --dalennod to decrypt")

	flag.BoolVar(&FlagVals.Import, "i", false, "Import bookmarks from a browser")
	flag.BoolVar(&FlagVals.Import, "import", false, "Import bookmarks from a browser")
	flag.StringVar(&FlagVals.Firefox, "firefox", "", "Import bookmarks from Firefox. Use alongside -i flag")
	flag.StringVar(&FlagVals.Chromium, "chromium", "", "Import bookmarks from Chromium. Use alongside -i flag")
	flag.StringVar(&FlagVals.Dalennod, "dalennod", "", "Import bookmarks exported Dalennod JSON. Use alongside -i flag")

	flag.BoolVar(&FlagVals.Where, "where", false, "Print config and logs directory location")

	flag.BoolVar(&FlagVals.RedoCompletion, "redo-completion", false, "Set CLI completion again")

	flag.BoolVar(&FlagVals.Profile, "profile", false, "Show profile names found in local directory")
	flag.StringVar(&FlagVals.Switch, "switch", "", "Switch profiles. Must use alongside --profile flag")

	flag.BoolVar(&FlagVals.FixDB, "fix-db", false, "Apply appropriate database fixes and updates")
}

func SetCompletion() {
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
	cachePath, err := os.UserCacheDir()
	if err != nil {
		log.Println("error finding home directory. ERROR:", err)
		return
	}

	fishLocalPath := filepath.Join(cachePath, "fish", "generated_completions")
	fishLocalStat, err := os.Stat(fishLocalPath)
	if err != nil {
		log.Println("error getting fish shell local directory info. ERROR", err)
		return
	}

	if !fishLocalStat.IsDir() {
		os.MkdirAll(fishLocalPath, 0755)
	}

	fishCompletionPath := filepath.Join(fishLocalPath, "dalennod.fish")
	fishCompletionFile, err := os.Create(fishCompletionPath)
	if err != nil {
		log.Println("error creating fish completion file. ERROR:", err)
	}
	defer fishCompletionFile.Close()

	sb := strings.Builder{}
	sb.WriteString("# Autogenerated from dalennod program\n")
	sb.WriteString("# $XDG_CACHE_HOME/fish/generated_completions/dalennod.fish\n\n")
	sb.WriteString("complete -c dalennod -s a -l add -d 'Add a bookmark entry to the database'\n")
	sb.WriteString("complete -c dalennod -s V -l view-all -d 'View all bookmarks'\n")
	sb.WriteString("complete -c dalennod -s v -l view -r -d 'View specific bookmark using its ID'\n")
	sb.WriteString("complete -c dalennod -s u -l update -r -d 'Update specific bookmark using its ID'\n")
	sb.WriteString("complete -c dalennod -s r -l remove -r -d 'Remove specific bookmark using its ID'\n")
	sb.WriteString("complete -c dalennod -s b -l backup -d 'Start backup process'\n")
	sb.WriteString("complete -c dalennod -s s -l serve -d 'Start webserver locally for Web UI & Extension'\n")
	sb.WriteString("complete -c dalennod -s h -l help -d 'Shows help message'\n")
	sb.WriteString("complete -c dalennod -l json -d 'Dump entire DB in JSON. Use alongside -b, --backup flag'\n")
	sb.WriteString("complete -c dalennod -l crypt -d 'Encrypt/decrypt the JSON backup'\n")
	sb.WriteString("complete -c dalennod -l where -d 'Print config and logs directory location'\n")
	sb.WriteString("complete -c dalennod -l redo-completion -d 'Set CLI completion again'\n")
	sb.WriteString("complete -c dalennod -l profile -d 'Show profile names found in local directory'\n")
	sb.WriteString("complete -c dalennod -l switch -d 'Switch profiles. Must use alongside --profile flag'\n\n")
	sb.WriteString("# import options\n")
	sb.WriteString("complete -c dalennod -s i -l import -d 'Import bookmarks from a browser'\n")
	sb.WriteString("complete -c dalennod -l firefox -d 'Import bookmarks from Firefox. Must use alongside -i, --import option'\n")
	sb.WriteString("complete -c dalennod -l chromium -d 'Import bookmarks from Chromium. Must use alongside -i, --import option'\n")
	sb.WriteString("complete -c dalennod -l dalennod -d 'Import bookmarks from exported Dalennod JSON. Must use alongside -i, --import option'\n")

	if _, err := fishCompletionFile.WriteString(sb.String()); err != nil {
		log.Println("error writing to fish completion file. ERROR:", err)
		return
	}
	fmt.Println("Created and set file for command line completion on fish shell.")
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

	sb := strings.Builder{}
	sb.WriteString("\n# Autogenerated from dalennod program\n\n")
	sb.WriteString("_dalennod() {\n")
	sb.WriteString("    local cur prev opts\n")
	sb.WriteString("    COMPREPLY=()\n")
	sb.WriteString("    cur=\"${COMP_WORDS[COMP_CWORD]}\"\n")
	sb.WriteString("    prev=\"${COMP_WORDS[COMP_CWORD-1]}\"\n")
	sb.WriteString("    opts=\"-s --serve -a --add -r --remove -u --update -v --view -V --view-all -i --import --firefox --chromium --dalennod -b --backup --json --crypt --where --profile --switch --redo-completion\"\n\n")
	sb.WriteString("    if [[ ${cur} == -* ]] ; then\n")
	sb.WriteString("        COMPREPLY=( $(compgen -W \"${opts}\" -- ${cur}) )\n")
	sb.WriteString("        return 0\n")
	sb.WriteString("    fi\n")
	sb.WriteString("}\n\n")
	sb.WriteString("complete -F _dalennod dalennod\n\n")
	sb.WriteString("# End of autogenerated lines from dalennod program\n")

	fmt.Printf("%s\nLines above can be auto-appended to .bashrc [at %s]\nto get command line completion for 'dalennod' in your shell environment.\n", sb.String(), bashrcPath)
	if !askUserConfirmation("Proceed?") {
		return
	}

	if _, err := bashrcFile.WriteString(sb.String()); err != nil {
		log.Fatalln("error writing to .bashrc. ERROR:", err)
	}
	fmt.Println("Lines appended. Reload .bashrc for command line completion.")
}

func zshCompletion() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Println("error finding home directory. ERROR:", err)
	}

	compdefPath := filepath.Join(constants.CONFIG_PATH, "zsh-completion")
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

	sb := strings.Builder{}
	sb.WriteString("#compdef dalennod\n")
	sb.WriteString("# Autogenerated from dalennod program\n\n")
	sb.WriteString("_arguments -s \\\n")
	sb.WriteString("    '(-s,--serve)'{-s,--serve}'[Start webserver locally for Web UI & Extension]' \\\n")
	sb.WriteString("    '(-a,--add)'{-a,--add}'[Add a bookmark entry to the database]' \\\n")
	sb.WriteString("    '(-r,--remove)'{-r,--remove}'[Remove specific bookmark using its ID]' \\\n")
	sb.WriteString("    '(-u,--update)'{-u,--update}'[Update specific bookmark using its ID]' \\\n")
	sb.WriteString("    '(-v,--view)'{-v,--view}'[View specific bookmark using its ID]' \\\n")
	sb.WriteString("    '(-V,--view-all)'{-V,--view-all}'[View all bookmarks]' \\\n")
	sb.WriteString("    '(-i,--import)'{-i,--import}'[Import bookmarks from a browser]' \\\n")
	sb.WriteString("    '(-h,--help)'{-h,--help}'[Shows help message]' \\\n")
	sb.WriteString("    '--firefox[Import bookmarks from Firefox. Must use alongside -i, --import option]' \\\n")
	sb.WriteString("    '--chromium[Import bookmarks from Chromium. Must use alongside -i, --import option]' \\\n")
	sb.WriteString("    '--dalennod[Import bookmarks exported Dalennod JSON. Must use alongside -i, --import option]' \\\n")
	sb.WriteString("    '(-b,--backup)'{-b,--backup}'[Start backup process]' \\\n")
	sb.WriteString("    '--json[Dump entire DB in JSON. Use alongside -b, --backup flag]' \\\n")
	sb.WriteString("    '--crypt[Encrypt or decrypt the JSON backup]' \\\n")
	sb.WriteString("    '--where[Print config and logs directory location]'\\\n")
	sb.WriteString("    '--redo-completion[Set CLI completion again]'\\\n")
	sb.WriteString("    '--profile[Show profile names found in local directory]' \\\n")
	sb.WriteString("    '--switch[Switch profiles. Must use alongside --profile flag]'\n")

	if _, err := compdefDalennodFile.WriteString(sb.String()); err != nil {
		log.Println("error writing to zsh's compdef _dalennod file. ERROR:", err)
		return
	}

	zshrcPath := filepath.Join(homePath, ".zshrc")
	zshrcFile, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("error opening zshrc. ERROR:", zshrcPath)
	}
	defer zshrcFile.Close()

	sb.Reset()
	sb.WriteString("\n# Autogenerated from dalennod program\n")
	sb.WriteString("export FPATH=" + compdefPath + ":$FPATH\n")
	sb.WriteString("autoload -Uz compinit && compinit\n")
	sb.WriteString("# End of autogenerated lines from dalennod program\n")

	fmt.Printf("%s\nLines above can be auto-appended to .zshrc [at %s]\nto get command line completion for 'dalennod' in your shell environment.\n", sb.String(), zshrcPath)
	if !askUserConfirmation("Proceed?") {
		return
	}

	if _, err := zshrcFile.WriteString(sb.String()); err != nil {
		log.Println("error writing into .zshrc. ERROR:", err)
	}
	fmt.Println("Lines appended. Reload .zshrc or shell for command line completion.")
}

func askUserConfirmation(ask string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/n]: ", ask)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("error reading input. ERROR:", err)
	}
	response = strings.ToLower(strings.TrimSpace(response))
	if response == "y" {
		return true
	} else if response == "n" {
		return false
	} else {
		fmt.Println("ERROR: unrecognized input. leaving")
		return false
	}
}

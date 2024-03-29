package setup

import (
	"encoding/json"
	"log"
	"os"
)

type CFG struct {
	FirstRun        bool `json:"firstRun"`
	DarkMode        bool `json:"darkMode"`
	AlwaysOverwrite bool `json:"alwaysOverwrite"`
}

const CFG_FILE string = "config.json"

func CfgSetup() {
	_, err := ReadCfg()
	if err != nil {
		log.Fatalln(err)
	} else {
		return
	}

	var config CFG = CFG{
		FirstRun:        false,
		DarkMode:        false,
		AlwaysOverwrite: true,
	}
	cfgJson, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}

	cfgDir, err := ConfigDir()
	if err != nil {
		log.Fatalln(err)
	}

	cfgFile, err := os.Create(cfgDir + CFG_FILE)
	if err != nil {
		log.Fatalln(err)
	}
	defer cfgFile.Close()

	cfgFile.Write(cfgJson)
}

func ReadCfg() (CFG, error) {
	var conf CFG

	cfgDir, err := ConfigDir()
	if err != nil {
		return conf, err
	}

	cfgContent, err := os.ReadFile(cfgDir + CFG_FILE)
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(cfgContent, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

func WriteCfg(fr, dm, ao bool) error {
	var config CFG = CFG{
		FirstRun:        fr,
		DarkMode:        dm,
		AlwaysOverwrite: ao,
	}
	cfgJson, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	cfgDir, err := ConfigDir()
	if err != nil {
		return err
	}

	cfgFile, err := os.Create(cfgDir + CFG_FILE)
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	cfgFile.Write(cfgJson)

	return nil
}

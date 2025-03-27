package setup

import (
    "dalennod/internal/constants"
    "encoding/json"
    "log"
    "os"
    "path/filepath"
)

type CFG struct {
    FirstRun bool   `json:"firstRun"`
    Host     string `json:"host"`
    Port     string `json:"port"`
}

func configSetup(cfgDir string) {
    config := CFG{
        FirstRun: true,
        Host: "",
        Port: constants.WEBUI_PORT,
    }
    cfgJson, err := json.MarshalIndent(config, "", "\t")
    if err != nil {
        log.Fatalln(err)
    }

    cfgFile, err := os.Create(filepath.Join(cfgDir, constants.CONFIG_FILENAME))
    if err != nil {
        log.Fatalln(err)
    }
    defer cfgFile.Close()

    if _, err = cfgFile.Write(cfgJson); err != nil {
        log.Println(err)
    }
}

func readCfg() (CFG, error) {
    var conf CFG

    cfgContent, err := os.ReadFile(filepath.Join(constants.CONFIG_PATH, constants.CONFIG_FILENAME))
    if err != nil {
        return conf, err
    }

    err = json.Unmarshal(cfgContent, &conf)
    if err != nil {
        return conf, err
    }

    return conf, nil
}

func writeCfg(firstRun bool, host string, port string) error {
    config := CFG{
        FirstRun: firstRun,
        Host: host,
        Port: port,
    }
    cfgJson, err := json.MarshalIndent(config, "", "\t")
    if err != nil {
        return err
    }

    cfgFile, err := os.Create(filepath.Join(constants.CONFIG_PATH, constants.CONFIG_FILENAME))
    if err != nil {
        return err
    }
    defer cfgFile.Close()

    cfgFile.Write(cfgJson)

    return nil
}

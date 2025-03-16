package setup

import (
    "dalennod/internal/constants"
    "encoding/json"
    "log"
    "os"
    "path/filepath"
)

type CFG struct {
    FirstRun bool `json:"firstRun"`
}

func configSetup(cfgDir string) {
    config := CFG{
        FirstRun: true,
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

func ReadCfg() (CFG, error) {
    var conf CFG

    cfgDir, err := ConfigDir()
    if err != nil {
        return conf, err
    }

    cfgContent, err := os.ReadFile(filepath.Join(cfgDir, constants.CONFIG_FILENAME))
    if err != nil {
        return conf, err
    }

    err = json.Unmarshal(cfgContent, &conf)
    if err != nil {
        return conf, err
    }

    return conf, nil
}

func WriteCfg(firstRun bool) error {
    config := CFG{
        FirstRun: firstRun,
    }
    cfgJson, err := json.MarshalIndent(config, "", "\t")
    if err != nil {
        return err
    }

    cfgDir, err := ConfigDir()
    if err != nil {
        return err
    }

    cfgFile, err := os.Create(filepath.Join(cfgDir, constants.CONFIG_FILENAME))
    if err != nil {
        return err
    }
    defer cfgFile.Close()

    cfgFile.Write(cfgJson)

    return nil
}

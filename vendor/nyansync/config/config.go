package config

import (
    "os"
    "log"
    "io/ioutil"
    "path/filepath"
    "gopkg.in/yaml.v2"

    "github.com/pkg/errors"
    "github.com/crgimenes/goconfig" // config
    _ "github.com/crgimenes/goconfig/yaml"

    "nyansync/location"
)

func Save(cfg *Config) error {
    if err := os.MkdirAll(goconfig.Path, 0755); err != nil {
        return errors.Wrap(err, "create config dir")
    }

    data, err := yaml.Marshal(cfg)
    if err != nil {
        return errors.Wrap(err, "yaml marshalling")
    }

    path := filepath.Join(goconfig.Path, goconfig.File)
    if err := ioutil.WriteFile(path, data, 0640); err != nil {
        return errors.Wrap(err, "create config file")
    }

    return nil
}

func Load() (*Config) {
    cfg := &Config{}

    var err error
    if goconfig.Path, err = location.DefaultConfigDir("NyanSync"); err != nil {
        log.Println("Unable to set default config path")
        panic(err)
    }
    goconfig.File = "nyansync.yaml"

    if err := goconfig.Parse(cfg); err != nil {
        log.Println("Unable to read config file")
        panic(err)
    }

    if err := Save(cfg); err != nil {
        log.Println("Unable to save the config file")
        panic(err)
    }

    return cfg
}

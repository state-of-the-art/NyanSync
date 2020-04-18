package config

import (
    "os"
    "log"
    "io/ioutil"
    "path/filepath"
    "gopkg.in/yaml.v2"

    "github.com/crgimenes/goconfig"        // config framework
    _ "github.com/crgimenes/goconfig/yaml"

    "nyansync/location"
)

func Save(cfg *Config) {
    if err := os.MkdirAll(goconfig.Path, 0755); err != nil {
        panic("Error create config dir")
    }

    data, err := yaml.Marshal(cfg)
    if err != nil {
        panic("Error during yaml marshalling")
    }

    path := filepath.Join(goconfig.Path, goconfig.File)
    if err := ioutil.WriteFile(path, data, 0640); err != nil {
        panic("Error during write config file")
    }
}

func Load() (*Config) {
    cfg := &Config{}

    goconfig.Path = location.DefaultConfigDir()
    goconfig.File = "nyansync.yaml"

    if err := goconfig.Parse(cfg); err != nil {
        log.Println("Unable to read config file")
        panic(err)
    }

    Save(cfg)

    return cfg
}

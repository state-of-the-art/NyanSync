package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/crgimenes/goconfig" // config framework
	_ "github.com/crgimenes/goconfig/yaml"
	"gopkg.in/yaml.v2"

	"github.com/state-of-the-art/NyanSync/lib/location"
)

func (cfg *Config) SaveNow() {
	log.Println("Saving yaml config", cfg.FilePathGet())
	if err := os.MkdirAll(goconfig.Path, 0755); err != nil {
		log.Panic("Error create config dir", err)
	}

	cfg.RLock()
	defer cfg.RUnlock()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Panic("Error during yaml marshalling", err)
	}

	if err := ioutil.WriteFile(cfg.FilePathGet(), data, 0640); err != nil {
		log.Panic("Error during write config file", err)
	}
}

func (cfg *Config) Save() {
	if !cfg.SaveLock() {
		return
	}

	go func() {
		// Wait for 5 seconds and save the config
		time.Sleep(5 * time.Second)
		defer cfg.SaveUnlock()
		cfg.SaveNow()
	}()
}

func Load() *Config {
	cfg := &Config{}

	goconfig.Path = location.DefaultConfigDir()
	goconfig.File = "nyansync.yaml"

	cfg.FilePathSet(filepath.Join(goconfig.Path, goconfig.File))
	if err := goconfig.Parse(cfg); err != nil {
		log.Panic("Unable to read config file", err)
	}

	cfg.Save()

	return cfg
}

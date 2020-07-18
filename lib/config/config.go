package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kkyr/fig" // config framework
	"gopkg.in/yaml.v2"

	"github.com/state-of-the-art/NyanSync/lib/location"
)

var ConfigName = "nyanshare.yaml"

func (cfg *Config) SaveNow() {
	log.Println("[INFO] Saving yaml config", cfg.FilePathGet())
	if err := os.MkdirAll(filepath.Dir(cfg.FilePathGet()), 0755); err != nil {
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

func Load() {
	cfg = &Config{}

	if err := fig.Load(cfg,
		fig.UseEnv("app"),
		fig.File(ConfigName),
		fig.Dirs(location.DefaultConfigDir()),
	); err != nil {
		log.Panic("Unable to read config file", err)
	}
	cfg.FilePathSet(filepath.Join(location.DefaultConfigDir(), ConfigName))

	if cfg.Sources == nil {
		cfg.Sources = make(map[string]Source)
	}
	if cfg.Receivers == nil {
		cfg.Receivers = make(map[string]Receiver)
	}

	cfg.Save()
}

// Core configuration
var cfg = &Config{}

func Cfg() *Config {
	return cfg
}

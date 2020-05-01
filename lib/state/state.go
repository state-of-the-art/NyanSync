package state

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/state-of-the-art/NyanSync/lib/config"
	"github.com/state-of-the-art/NyanSync/lib/crypt"
	"github.com/state-of-the-art/NyanSync/lib/location"
)

const (
	init_admin_login         = "admin"
	init_admin_password_file = "admin_init_pass.txt"
)

func Init(config_state_path string) {
	state.FilePathSet(location.RealFilePath(config_state_path))

	if state.Load() {
		if state.Sources == nil {
			state.Sources = make(map[string]Source)
		}
		// Do not need to init the admin user if state is loaded
		return
	}
	// TODO: Remove duplication
	if state.Sources == nil {
		state.Sources = make(map[string]Source)
	}

	// Create admin password, create admin user and store password as admin file
	admin_pass := crypt.RandString(32)
	user := UserGet(init_admin_login)
	user.Set(admin_pass, "Administrator", true)

	state.SaveNow()

	admin_password_file := filepath.Join(filepath.Dir(state.FilePathGet()), init_admin_password_file)
	if err := ioutil.WriteFile(admin_password_file, []byte(admin_pass), 0400); err != nil {
		log.Panic("Unable to write admin password file", err)
	}
}

func SourcesUpdate(cfg_sources map[string]config.Source) {
	// Update the state sources according the config
	for id := range cfg_sources {
		s := SourceGet(id)
		s.Set(cfg_sources[id].Url, cfg_sources[id].Type)
	}

	// Removing not existing sources from state
	if len(cfg_sources) != len(state.Sources) {
		for id := range state.Sources {
			if _, ok := cfg_sources[id]; !ok {
				SourceRemove(id)
			}
		}
	}

	state.Save()
}

// Struct to store state
type st struct {
	config.Base

	Users   []User
	Sources map[string]Source
}

// Current application state
var state = &st{}

func (cfg *st) SaveNow() {
	log.Println("Saving json", cfg.FilePathGet())
	if err := os.MkdirAll(filepath.Dir(cfg.FilePathGet()), 0750); err != nil {
		log.Panic("Unable to create save dir", err)
	}

	cfg.RLock()
	defer cfg.RUnlock()

	data, err := json.Marshal(cfg)
	if err != nil {
		log.Panic("Error during serializing of the base config struct", err)
	}

	if err := ioutil.WriteFile(cfg.FilePathGet(), data, 0640); err != nil {
		log.Panic("Error during write base config to file", err)
	}
}

func (cfg *st) Save() {
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

func (cfg *st) Load() bool {
	log.Println("[DEBUG] Loading from file", cfg.FilePathGet())
	if file, err := os.OpenFile(cfg.FilePathGet(), os.O_RDONLY, 640); err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		cfg.Lock()
		defer cfg.Unlock()
		if err = decoder.Decode(cfg); err != nil {
			log.Panic("Unable to decode loaded file", err)
		}
		return true
	}
	return false
}

func SourcesList() []Source {
	// TODO: actual sources
	out := make([]Source, 0, len(state.Sources))

	for _, value := range state.Sources {
		out = append(out, value)
	}
	return out
}

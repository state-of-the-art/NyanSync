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

func Init() {
	state.FilePathSet(location.RealFilePath(config.Cfg().StateFilePath))

	if !state.Load() {
		// Create admin password, user and store password as admin file
		admin_pass := crypt.RandString(32)
		user := UserGet(init_admin_login)
		user.Set(admin_pass, "Administrator", true)

		if err := os.MkdirAll(filepath.Dir(state.FilePathGet()), 0750); err != nil {
			log.Panic("Unable to create dir: ", err)
		}
		admin_password_file := filepath.Join(filepath.Dir(state.FilePathGet()), init_admin_password_file)
		if err := ioutil.WriteFile(admin_password_file, []byte(admin_pass), 0400); err != nil {
			log.Panic("Unable to write admin password file: ", err)
		}
	}

	if state.Sources == nil {
		state.Sources = make(map[string]Source)
	}

	SourcesUpdateFromConfig()
}

func SourcesUpdateFromConfig() {
	cfg_sources := config.Cfg().Sources
	// Update the state sources according the config
	for id := range cfg_sources {
		s := SourceGet(id)
		s.Set(cfg_sources[id].Uri, cfg_sources[id].Type)
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

func (s *st) SaveNow() {
	log.Println("Saving json", s.FilePathGet())
	if err := os.MkdirAll(filepath.Dir(s.FilePathGet()), 0750); err != nil {
		log.Panic("Unable to create save dir: ", err)
	}

	s.RLock()
	defer s.RUnlock()

	data, err := json.Marshal(s)
	if err != nil {
		log.Panic("Error during serializing of the base config struct", err)
	}

	if err := ioutil.WriteFile(s.FilePathGet(), data, 0640); err != nil {
		log.Panic("Error during write base config to file", err)
	}
}

func (s *st) Save() {
	if !s.SaveLock() {
		return
	}

	go func() {
		// Wait for 5 seconds and save the config
		time.Sleep(5 * time.Second)
		defer s.SaveUnlock()
		s.SaveNow()
	}()
}

func (s *st) Load() bool {
	log.Println("[DEBUG] Loading from file", s.FilePathGet())
	if file, err := os.OpenFile(s.FilePathGet(), os.O_RDONLY, 640); err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		s.Lock()
		defer s.Unlock()
		if err = decoder.Decode(s); err != nil {
			log.Panic("Unable to decode loaded file", err)
		}
		return true
	}
	return false
}

func UsersList() []User {
	return state.Users
}

func SourcesList() []Source {
	out := make([]Source, 0, len(state.Sources))

	for _, value := range state.Sources {
		out = append(out, value)
	}
	return out
}

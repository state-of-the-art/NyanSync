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
		user.PasswordSet(admin_pass)
		user.Roles = []string{"admin"}
		user.Name = "Administrator"
		user.Manager = "" // Noone manages admin
		user.Save()

		if err := os.MkdirAll(filepath.Dir(state.FilePathGet()), 0750); err != nil {
			log.Panic("Unable to create dir: ", err)
		}
		admin_password_file := filepath.Join(filepath.Dir(state.FilePathGet()), init_admin_password_file)
		if err := ioutil.WriteFile(admin_password_file, []byte(admin_pass), 0400); err != nil {
			log.Panic("Unable to write admin password file: ", err)
		}
	}

	SourcesUpdateFromConfig()
}

func SourcesUpdateFromConfig() {
	cfg_sources := config.Cfg().Sources
	// Update the state sources according the config
	for id := range cfg_sources {
		s := SourceGet(id)
		s.Set(s.Manager, cfg_sources[id].Uri, cfg_sources[id].Type)
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

	Users   map[string]User
	Sources map[string]Source
	state_w bool

	Access   map[string]Access `json:"-"` // Stored to another file
	access_w bool
}

// Current application state
var state = &st{}

func (s *st) SaveNow() {
	if s.state_w {
		s.RLock()
		saveJson(s, s.FilePathGet())
		s.RUnlock()

		s.Lock()
		s.state_w = false
		s.Unlock()
	}
	if s.access_w {
		s.RLock()
		saveJson(s.Access, location.RealFilePath(config.Cfg().AccessFilePath))
		s.RUnlock()

		s.Lock()
		s.access_w = false
		s.Unlock()
	}
}

func saveJson(s interface{}, path string) {
	log.Println("[DEBUG] Saving json", path)
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		log.Panic("Unable to create save dir: ", err)
	}

	data, err := json.Marshal(s)
	if err != nil {
		log.Panic("Error during serializing for the file ", path, ": ", err)
	}

	if err := ioutil.WriteFile(path, data, 0640); err != nil {
		log.Panic("Error during write to file ", path, ": ", err)
	}
}

func (s *st) Save() {
	go s.save(&s.state_w)
}

func (s *st) SaveAccess() {
	go s.save(&s.access_w)
}

func (s *st) save(flag *bool) {
	s.Lock()
	defer s.Unlock()
	*flag = true

	if !s.SaveLock() {
		return
	}

	go func() {
		defer s.SaveUnlock()
		// Wait for 5 seconds and save the config
		time.Sleep(5 * time.Second)
		s.SaveNow()
	}()
}

func (s *st) Load() bool {
	var ok = s.loadState()
	s.loadAccess()

	if state.Users == nil {
		state.Users = make(map[string]User)
	}
	if state.Sources == nil {
		state.Sources = make(map[string]Source)
	}
	if state.Access == nil {
		state.Access = make(map[string]Access)
	}

	return ok
}

func (s *st) loadState() bool {
	log.Println("[DEBUG] Loading state from file", s.FilePathGet())
	if file, err := os.OpenFile(s.FilePathGet(), os.O_RDONLY, 640); err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		s.Lock()
		defer s.Unlock()
		if err = decoder.Decode(s); err != nil {
			log.Panic("Unable to decode loaded file: ", err)
		}
		return true
	}
	return false
}

func (s *st) loadAccess() {
	path := location.RealFilePath(config.Cfg().AccessFilePath)
	log.Println("[DEBUG] Loading access from file", path)
	if file, err := os.OpenFile(path, os.O_RDONLY, 640); err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		s.Lock()
		defer s.Unlock()
		if err = decoder.Decode(&s.Access); err != nil {
			log.Panic("Unable to decode loaded file: ", err)
		}
	}
}

func UsersList() []User {
	out := make([]User, 0, len(state.Users))

	state.RLock()
	defer state.RUnlock()

	for _, value := range state.Users {
		out = append(out, value)
	}
	return out
}

func SourceList() []Source {
	out := make([]Source, 0, len(state.Sources))

	state.RLock()
	defer state.RUnlock()

	for _, value := range state.Sources {
		out = append(out, value)
	}
	return out
}

func AccessList() []Access {
	out := make([]Access, 0, len(state.Access))

	state.RLock()
	defer state.RUnlock()

	for _, value := range state.Access {
		out = append(out, value)
	}
	return out
}

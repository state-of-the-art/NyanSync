package state

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/euroteltr/rbac"

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
		user.Set(admin_pass, "Administrator", "", true)

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

	Users   []User
	Sources map[string]Source
	state_w bool

	Access   map[string]Access `json:"-"` // Stored to another file
	access_w bool

	RBAC   *rbac.RBAC `json:"-"` // Stored to another file
	rbac_w bool
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
	if s.rbac_w {
		s.RLock()
		saveJson(s.RBAC, location.RealFilePath(config.Cfg().RBACFilePath))
		s.RUnlock()

		s.Lock()
		s.rbac_w = false
		s.Unlock()
	}
}

func saveJson(s interface{}, path string) {
	log.Println("Saving json", path)
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

func (s *st) SaveRBAC() {
	go s.save(&s.rbac_w)
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
	s.loadRBAC()

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

func (s *st) loadRBAC() {
	s.RBAC = rbac.New(rbac.NewConsoleLogger())
	path := location.RealFilePath(config.Cfg().RBACFilePath)
	log.Println("[DEBUG] Loading rbac from file", path)
	if file, err := os.OpenFile(path, os.O_RDONLY, 640); err == nil {
		defer file.Close()
		if err = s.RBAC.LoadJSON(file); err != nil {
			log.Panic("Unable to load rbac file: ", err)
		}
	} else {
		// Init default RBAC rules
		user, _ := state.RBAC.RegisterPermission("user", "User resource", rbac.Create, rbac.Read, rbac.Update, rbac.Delete)
		access, _ := state.RBAC.RegisterPermission("access", "Access resource", rbac.Create, rbac.Read, rbac.Update, rbac.Delete)
		source, _ := state.RBAC.RegisterPermission("source", "Source resource", rbac.Create, rbac.Read, rbac.Update, rbac.Delete)

		ApproveAction := rbac.Action("approve")

		admin, _ := state.RBAC.RegisterRole("admin", "Admin role")

		state.RBAC.Permit(admin.ID, user, rbac.CRUD, ApproveAction)
		state.RBAC.Permit(admin.ID, access, rbac.CRUD, ApproveAction)
		state.RBAC.Permit(admin.ID, source, rbac.CRUD, ApproveAction)

		state.SaveRBAC()
	}
}

func UsersList() []User {
	return state.Users
}

func SourceList() []Source {
	out := make([]Source, 0, len(state.Sources))

	for _, value := range state.Sources {
		out = append(out, value)
	}
	return out
}

func AccessList() []Access {
	out := make([]Access, 0, len(state.Access))

	for _, value := range state.Access {
		out = append(out, value)
	}
	return out
}

package state

import (
	"log"
	"os"
	"sync"
	"io/ioutil"
	"encoding/json"
	"path/filepath"

	"github.com/state-of-the-art/NyanSync/lib/config"
	"github.com/state-of-the-art/NyanSync/lib/crypt"
	"github.com/state-of-the-art/NyanSync/lib/location"
)

const (
	init_admin_login         = "admin"
	init_admin_password_file = "admin_init_pass.txt"
)

func Init(config_state_path string) {
	state.file_path = location.RealFilePath(config_state_path)

	log.Println("[DEBUG] Opening stored state", state.file_path)
	if file, err := os.OpenFile(state.file_path, os.O_RDONLY, 640); err == nil {
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(state); err != nil {
			panic(err)
		}
		return
	} else {
		log.Println("[INFO] Unable to open state - creating a new one", err)
	}

	parent := filepath.Dir(state.file_path)
	if err := os.MkdirAll(parent, 0750); err != nil {
		panic("Unable to create access file dir")
	}

	// Create admin password, create admin user and store password as admin file
	admin_pass := crypt.RandString(32)
	user := UserGet(init_admin_login)
	user.Set(admin_pass, "Administrator", true)

	if err := ioutil.WriteFile(filepath.Join(parent, init_admin_password_file), []byte(admin_pass), 0400); err != nil {
		panic("Unable to write admin password file")
	}

	Save()
}

func SourcesUpdate(cfg_sources []config.Source) {
}

func Save() {
	// Save state file
	data, err := json.Marshal(state)
	if err != nil {
		panic("Error during serializing of state struct")
	}

	if err := ioutil.WriteFile(state.file_path, data, 0640); err != nil {
		panic("Error during write the state file")
	}
}

// Struct to store state
type st struct {
	sync.RWMutex
	file_path string

	Users   []User
	Sources []Source
}

// Current application state
var state = &st{}

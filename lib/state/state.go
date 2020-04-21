package state

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/state-of-the-art/NyanSync/lib/crypt"
	"github.com/state-of-the-art/NyanSync/lib/location"
)

const (
	InitAdminLogin        = "admin"
	InitAdminPasswordFile = "admin_init_pass.txt"
)

func Init(config_state_path string) {
	state.file_path = location.RealFilePath(config_state_path)

	if file, err := os.Open(state.file_path); os.IsExist(err) {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(state)
		return
	}

	// No access file found, so we need to create one

	parent := filepath.Dir(state.file_path)
	if err := os.MkdirAll(parent, 0750); err != nil {
		panic("Unable to create access file dir")
	}

	// Create admin password, create admin user and store password as admin file
	admin_pass := crypt.RandString(32)
	UserSet(InitAdminLogin, admin_pass, "Administrator", true)

	if err := ioutil.WriteFile(filepath.Join(parent, InitAdminPasswordFile), []byte(admin_pass), 0400); err != nil {
		panic("Unable to write admin password file")
	}

	Save()
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

	Users []User
}

// Current application state
var state = &st{}

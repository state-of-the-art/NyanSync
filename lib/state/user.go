package state

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/state-of-the-art/NyanSync/lib/crypt"
)

const (
	InvalidUserManager = "Invalid user Manager"
)

type User struct {
	sync.RWMutex

	Login    string     // Login id of the user
	Name     string     // Real name of the user
	Manager  string     // Controlling manager
	Roles    []string   // List of the RBAC roles
	Init     bool       // First time created user - should be recreated
	PassHash crypt.Hash // Hash + salt for the user password
}

func UserExists(login string) bool {
	state.RLock()
	defer state.RUnlock()
	_, ok := state.Users[login]
	return ok
}

func UserGet(login string) User {
	state.RLock()
	defer state.RUnlock()

	user, ok := state.Users[login]
	if !ok {
		user = User{
			Login: login,
			Init:  true,
		}
	}

	return user
}

func (u *User) Save() error {
	// Check manager is good
	if err := u.IsValid(); err != nil {
		return err
	}

	// Save current user to the array
	state.Lock()
	defer state.Unlock()
	state.Users[u.Login] = *u
	state.Save()

	return nil
}

func UserRemove(login string) {
	state.Lock()
	defer state.Unlock()
	delete(state.Users, login)
	state.Save()
}

func (u *User) IsValid() error {
	u.RLock()
	defer u.RUnlock()

	// Check manager is good
	if u.Manager != "" && !UserExists(u.Manager) {
		return errors.New(InvalidUserManager)
	}

	return nil
}

func (u *User) PasswordSet(password string) {
	u.Lock()
	defer u.Unlock()
	u.PassHash = crypt.Generate(password, nil)
}

func (u *User) PasswordCheck(password string) bool {
	u.RLock()
	defer u.RUnlock()
	return u.PassHash.IsEqual(password)
}

package state

import (
	"sync"

	"github.com/state-of-the-art/NyanSync/lib/crypt"
)

type User struct {
	sync.RWMutex

	Name     string     // Real name of the user
	Login    string     // Login id of the user
	Init     bool       // First time created user - should be recreated
	PassHash crypt.Hash // Hash + salt for the user password
}

func UserFind(login string) *User {
	state.Lock()
	defer state.Unlock()
	for _, user := range state.Users {
		if user.Login == login {
			return &user
		}
	}
	return nil
}

func UserGet(login string) *User {
	if puser := UserFind(login); puser != nil {
		return puser
	}

	state.Lock()
	defer state.Unlock()
	state.Users = append(state.Users, User{Login: login, Init: true})
	puser := &state.Users[len(state.Users)-1]
	return puser
}

func (u *User) Set(password string, name string, init bool) {
	u.Lock()
	defer u.Unlock()
	u.Name = name
	u.PassHash = crypt.Generate(password, nil)
	u.Init = init
}

func (u *User) CheckPassword(password string) bool {
	u.Lock()
	defer u.Unlock()
	return u.PassHash.IsEqual(password)
}

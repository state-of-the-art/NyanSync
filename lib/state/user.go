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

func UserGet(login string) *User {
	state.Lock()
	for _, user := range state.Users {
		if user.Login == login {
			return &user
		}
	}
	state.Users = append(state.Users, User{Login: login, Init: true})
	puser := &state.Users[len(state.Users)-1]
	state.Unlock()
	return puser
}

func UserSet(login string, password string, name string, init bool) {
	puser := UserGet(login)
	puser.Lock()
	puser.Name = name
	puser.PassHash = crypt.Generate(password, nil)
	puser.Init = init
	puser.Unlock()
}

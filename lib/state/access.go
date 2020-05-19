package state

import (
	"sync"

	"github.com/lithammer/shortuuid/v3"
	"github.com/pkg/errors"
)

const (
	InvalidAccessSourceId  = "Invalid access Source Id"
	InvalidAccessUserLogin = "Invalid access User Login"
	InvalidAccessManager   = "Invalid access Manager"
)

type Access struct {
	sync.RWMutex

	Id       string
	SourceId string
	Manager  string
	Path     string
	Users    []string
}

func AccessExists(id string) bool {
	_, ok := state.Access[id]
	return ok
}

func AccessGet(id string) Access {
	state.RLock()
	defer state.RUnlock()
	if o, ok := state.Access[id]; ok {
		return o
	}
	return Access{Id: id}
}

func AccessRemove(id string) {
	state.Lock()
	defer state.Unlock()
	delete(state.Access, id)
	state.SaveAccess()
}

func (a *Access) Set(source_id string, manager string, path string, users []string) {
	a.SourceId = source_id
	a.Manager = manager
	a.Path = path
	a.Users = users
	a.Save()
}

func (a *Access) Save() error {
	if err := a.IsValid(); err != nil {
		return err
	}
	state.Lock()
	defer state.Unlock()
	state.Access[a.Id] = *a
	state.SaveAccess()
	return nil
}

func (a *Access) IsValid() error {
	if !SourceExists(a.SourceId) {
		return errors.New(InvalidAccessSourceId)
	}
	for _, u := range a.Users {
		if UserFind(u) == nil {
			return errors.New(InvalidAccessUserLogin)
		}
	}
	if a.Manager != "" && UserFind(a.Manager) == nil {
		return errors.New(InvalidAccessManager)
	}
	return nil
}

func (a *Access) NewId() {
	a.Id = shortuuid.New()
	// Rarely the generated uuid could be already here - so need to check
	if _, exists := state.Access[a.Id]; exists {
		a.NewId()
	}
}

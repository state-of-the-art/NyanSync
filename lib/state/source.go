package state

import (
	"net/url"
	"sync"

	"github.com/pkg/errors"

	"github.com/state-of-the-art/NyanSync/lib/config"
)

const (
	InvalidSourceId      = "Invalid source Id"
	InvalidSourceUri     = "Invalid source URI"
	InvalidSourceType    = "Invalid source Type"
	InvalidSourceManager = "Invalid source Manager"
)

type Source struct {
	*config.Source
	sync.RWMutex

	Id      string
	Manager string
}

func SourceExists(id string) bool {
	_, ok := state.Sources[id]
	return ok
}

func SourceGet(id string) Source {
	state.RLock()
	defer state.RUnlock()
	if _, ok := state.Sources[id]; ok {
		return state.Sources[id]
	}
	return Source{&config.Source{}, sync.RWMutex{}, id, ""}
}

func SourceRemove(id string) {
	state.Lock()
	defer state.Unlock()
	delete(state.Sources, id)
	state.Save()
}

func (s *Source) Set(manager string, uri string, _type string /* TODO: Options */) {
	s.Manager = manager
	s.Uri = uri
	s.Type = _type
	s.Save()
}

func (s *Source) SaveRename(old_id string) error {
	if err := isValidId(old_id); err != nil {
		return err
	}
	if err := s.IsValid(); err != nil {
		return err
	}
	state.Lock()
	defer state.Unlock()
	if old_id != s.Id {
		delete(state.Sources, old_id)
	}
	state.Sources[s.Id] = *s
	state.Save()
	config.Cfg().SourceSet(s.Id, s.Source)
	return nil
}

func (s *Source) Save() error {
	if err := s.IsValid(); err != nil {
		return err
	}
	state.Lock()
	defer state.Unlock()
	state.Sources[s.Id] = *s
	state.Save()
	config.Cfg().SourceSet(s.Id, s.Source)
	return nil
}

func (s *Source) IsValid() error {
	if err := isValidId(s.Id); err != nil {
		return err
	}
	if _, err := url.ParseRequestURI(s.Uri); err != nil {
		return errors.Wrap(err, InvalidSourceUri)
	}
	if s.Type == "" {
		return errors.New(InvalidSourceType)
	}
	if s.Manager != "" && !UserExists(s.Manager) {
		return errors.New(InvalidSourceManager)
	}

	return nil
}

func isValidId(id string) error {
	if id == "" {
		return errors.New(InvalidSourceId)
	}
	return nil
}

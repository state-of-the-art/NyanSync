package state

import (
	"net/url"
	"sync"

	"github.com/pkg/errors"

	"github.com/state-of-the-art/NyanSync/lib/config"
)

const (
	InvalidSourceId   = "Invalid source Id"
	InvalidSourceUrl  = "Invalid source URL"
	InvalidSourceType = "Invalid source Type"
)

type Source struct {
	*config.Source
	sync.RWMutex

	Id string
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
	return Source{&config.Source{}, sync.RWMutex{}, id}
}

func SourceRemove(id string) {
	state.Lock()
	defer state.Unlock()
	delete(state.Sources, id)
	state.Save()
}

func (s *Source) Set(url string, _type string /* TODO: Options */) {
	s.Url = url
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
	return nil
}

func (s *Source) IsValid() error {
	if err := isValidId(s.Id); err != nil {
		return err
	}
	if _, err := url.ParseRequestURI(s.Url); err != nil {
		return errors.Wrap(err, InvalidSourceUrl)
	}
	if s.Type == "" {
		return errors.New(InvalidSourceType)
	}

	return nil
}

func isValidId(id string) error {
	if id == "" {
		return errors.New(InvalidSourceId)
	}
	return nil
}

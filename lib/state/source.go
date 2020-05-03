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
}

func (s *Source) Set(url string, _type string /* TODO: Options */) {
	s.Url = url
	s.Type = _type
	s.Save()
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
	if s.Id == "" {
		return errors.New(InvalidSourceId)
	}
	if _, err := url.ParseRequestURI(s.Url); err != nil {
		return errors.Wrap(err, InvalidSourceUrl)
	}
	if s.Type == "" {
		return errors.New(InvalidSourceType)
	}

	return nil
}

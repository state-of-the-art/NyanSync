package state

import (
	"sync"

	"github.com/state-of-the-art/NyanSync/lib/config"
)

type Source struct {
	*config.Source
	sync.RWMutex
}

func SourceFind(id string) *Source {
	state.Lock()
	defer state.Unlock()
	for _, source := range state.Sources {
		if source.Id == id {
			return &source
		}
	}
	return nil
}

func SourceGet(id string) *Source {
	if p := SourceFind(id); p != nil {
		return p
	}

	state.Lock()
	defer state.Unlock()
	state.Sources = append(state.Sources, Source{})
	p := &state.Sources[len(state.Sources)-1]
	p.Id = id
	return p
}

func (s *Source) Set(url string, _type string /* Options todo */) {
	s.Lock()
	defer s.Unlock()
	s.Url = url
	s.Type = _type
}

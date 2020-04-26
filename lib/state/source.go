package state

import (
	"sync"

	"github.com/state-of-the-art/NyanSync/lib/config"
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
	return Source{Id: id}
}

func SourceRemove(id string) {
	state.Lock()
	defer state.Unlock()
	delete(state.Sources, id)
}

func (s *Source) Set(url string, _type string /* Options todo */) {
	s.Url = url
	s.Type = _type
	state.Lock()
	defer state.Unlock()
	state.Sources[s.Id] = *s
}

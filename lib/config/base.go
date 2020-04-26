package config

import (
	"sync"
)

type Base struct {
	sync.RWMutex

	file_path string // Where to store the config

	// Delayed save lock
	save_lock   sync.Mutex
	save_active bool
}

func (base *Base) FilePathGet() string {
	return base.file_path
}

func (base *Base) FilePathSet(path string) {
	base.file_path = path
}

func (base *Base) SaveLock() bool {
	base.save_lock.Lock()
	defer base.save_lock.Unlock()
	if base.save_active {
		return false
	}
	base.save_active = true
	return true
}

func (base *Base) SaveUnlock() {
	base.save_lock.Lock()
	defer base.save_lock.Unlock()
	base.save_active = false
}

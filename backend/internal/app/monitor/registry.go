package monitor

import (
	"sync"

	"github.com/google/uuid"
)

var (
	activePollers = make(map[string]struct{}) // transactionID → active
	mu            sync.Mutex
)

func registerPoller(id uuid.UUID) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := activePollers[id.String()]; exists {
		return false // already running
	}

	activePollers[id.String()] = struct{}{}

	return true
}

func unregisterPoller(id uuid.UUID) {
	mu.Lock()
	defer mu.Unlock()

	delete(activePollers, id.String())
}

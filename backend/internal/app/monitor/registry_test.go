package monitor

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterPoller(t *testing.T) {
	id := uuid.New()

	t.Run("Register successful", func(t *testing.T) {
		ok := registerPoller(id)
		assert.True(t, ok, "First registration should succeed")

		defer unregisterPoller(id)
	})

	t.Run("Register duplicate fails", func(t *testing.T) {
		registerPoller(id)
		defer unregisterPoller(id)

		ok := registerPoller(id)
		assert.False(t, ok, "Double registration should fail")
	})

	t.Run("Unregister allows re-registration", func(t *testing.T) {
		registerPoller(id)
		unregisterPoller(id)

		// Nach dem Abmelden muss es wieder gehen
		ok := registerPoller(id)
		assert.True(t, ok, "Registration after unregister should work")

		unregisterPoller(id)
	})
}

func TestRegisterPolle_Concurrency(t *testing.T) {
	const goroutines = 100

	id := uuid.New()

	var wg sync.WaitGroup

	results := make(chan bool, goroutines)

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			results <- registerPoller(id)
		}()
	}

	wg.Wait()
	close(results)

	trueCount := 0

	for res := range results {
		if res {
			trueCount++
		}
	}

	assert.Equal(t, 1, trueCount, "Just one goroutine should succeed in registering the poller")

	// Cleanup
	unregisterPoller(id)
}

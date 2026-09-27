package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// errInjectedCreate is the sentinel a seed create hook returns to fail one Create.
var errInjectedCreate = errors.New("injected create failure")

// unknownCallSite is reported when the stack walk cannot find the seed.go frame,
// which would let the test cover nothing while still passing.
const unknownCallSite = "unknown"

// seedCreateHook observes every Create the seed performs, in order. Returning an
// error injects that failure into the Create.
type seedCreateHook func(index int, callSite string) error

// seedCallSite reports the seed.go "file:line" of the Create currently in flight.
// The seed frame sits a few levels above the callback GORM invokes, so walk outwards
// until it is found rather than hardcoding a skip count.
func seedCallSite() string {
	for skip := 1; skip < 10; skip++ {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		if filepath.Base(file) == "seed.go" {
			return fmt.Sprintf("seed.go:%d", line)
		}
	}

	return unknownCallSite
}

// runSeedWithCreateHook seeds a fresh database, lets hook observe and optionally fail
// each Create, and returns the error the seed reported.
//
// Every pass needs its own database: the seed always inserts the fixed guest code
// ABCDEFGHI and models.Guest.Code is unique, so a shared database would fail on that
// constraint instead of the injected one and pass for the wrong reason.
func runSeedWithCreateHook(t *testing.T, hook seedCreateHook) error {
	t.Helper()

	t.Chdir(t.TempDir())

	db, err := ConnectToDatabase("seedhook")
	require.NoError(t, err)

	t.Cleanup(func() { _ = CloseDatabase(db) })

	var injected error

	index := 0

	err = db.Callback().
		Create().
		Before("gorm:create").
		Register("test:seed_create_hook", func(tx *gorm.DB) {
			index++

			if hookErr := hook(index, seedCallSite()); hookErr != nil {
				injected = hookErr
				_ = tx.AddError(hookErr)
			}
		})
	require.NoError(t, err)

	seedErr := SeedDatabase(db, true)

	// A swallowed injection would let the subtests below pass without ever reaching
	// an error branch, so state that expectation here.
	if injected != nil {
		require.Error(t, seedErr, "seed swallowed the failure injected at create #%d", index)
	}

	return seedErr
}

// The count guards short-circuit before touching the database. SeedDatabase always
// passes a positive count, so call the helpers directly to keep that path covered.
func TestSeedHelpersAcceptZeroCounts(t *testing.T) {
	t.Chdir(t.TempDir())

	db, err := ConnectToDatabase("seedzero")
	require.NoError(t, err)

	t.Cleanup(func() { _ = CloseDatabase(db) })

	seed := NewDatabaseSeed(db)

	require.NoError(t, seed.seedUserGuests(0, 0, 0))
	require.NoError(t, seed.seedPurchases(0))
}

// Every seeding step has to report its failure, otherwise `database reset` reports
// success over a half-written fixture database. Those error branches are only
// reachable by making a specific Create fail, so walk each call site in turn: one
// recording pass notes the first create index per call site, then a pass per call site
// fails exactly that index and requires the seed to surface it.
//
// Seed calls gofakeit.Seed(1) before creating anything, so the create sequence is
// identical on every run and a recorded index keeps pointing at the same call site.
func TestSeedDatabaseReportsEveryCreateFailure(t *testing.T) {
	indexOfCallSite := map[string]int{}

	var callSites []string

	err := runSeedWithCreateHook(t, func(index int, callSite string) error {
		if _, seen := indexOfCallSite[callSite]; !seen {
			indexOfCallSite[callSite] = index
			callSites = append(callSites, callSite)
		}

		return nil
	})
	require.NoError(t, err)

	// If the stack walk stopped matching seed.go, every create would collapse into a
	// single "unknown" entry and the subtests below would assert nothing.
	require.NotContains(t, callSites, unknownCallSite, "create hook lost the seed.go call site")
	require.NotEmpty(t, callSites)

	for _, callSite := range callSites {
		t.Run(callSite, func(t *testing.T) {
			target := indexOfCallSite[callSite]

			err := runSeedWithCreateHook(t, func(index int, site string) error {
				if site == callSite && index == target {
					return errInjectedCreate
				}

				return nil
			})

			require.ErrorIs(t, err, errInjectedCreate, "seed must report the failure from %s", callSite)
		})
	}
}

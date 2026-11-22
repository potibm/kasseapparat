package sqlite

import (
	"fmt"

	"github.com/potibm/kasseapparat/internal/app/models"
)

// GuestlistStats sums up the entries of one guest list.
type GuestlistStats struct {
	// Entries is the number of list entries.
	Entries int `json:"entries"`
	// ExpectedGuests counts every entry with its additional guests.
	ExpectedGuests int `json:"expectedGuests"`
	// ArrivedGuests counts the people already checked in.
	ArrivedGuests int `json:"arrivedGuests"`
}

// GetGuestlistStats returns the stats of the given guest lists. A list without
// entries is missing from the map.
func (repo *Repository) GetGuestlistStats(ids []int) (map[int]GuestlistStats, error) {
	stats := make(map[int]GuestlistStats, len(ids))

	if len(ids) == 0 {
		return stats, nil
	}

	var rows []struct {
		GuestlistID    int
		Entries        int
		ExpectedGuests int
		ArrivedGuests  int
	}

	err := repo.db.Model(&models.Guest{}).
		Select("guestlist_id, COUNT(*) AS entries, "+
			"COALESCE(SUM(1 + additional_guests), 0) AS expected_guests, "+
			"COALESCE(SUM(attended_guests), 0) AS arrived_guests").
		Where("guestlist_id IN ?", ids).
		Group("guestlist_id").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count the guests of the lists: %w", err)
	}

	for _, row := range rows {
		stats[row.GuestlistID] = GuestlistStats{
			Entries:        row.Entries,
			ExpectedGuests: row.ExpectedGuests,
			ArrivedGuests:  row.ArrivedGuests,
		}
	}

	return stats, nil
}

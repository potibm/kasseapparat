package models

import (
	"sort"
	"strings"
	"time"
)

// List represents a guestlist
type ListEntry struct {
	GormOwnedModel
	GuestlistID          uint       `json:"guestlistId"`
	Guestlist            Guestlist  `json:"guestlist"`
	Name                 string     `json:"name" `
	Code                 *string    `json:"code" gorm:"unique"`
	AdditionalGuests     uint       `json:"additionalGuests" gorm:"default:0"`
	AttendedGuests       uint       `json:"attendedGuests" gorm:"default:0"`
	ArrivedAt            *time.Time `json:"arrivedAt"`
	ArrivalNote          *string    `json:"arrivalNote"`
	NotifyOnArrivalEmail *string    `json:"notifyOnArrivalEmail"`
	PurchaseID           *uint      `json:"purchaseId"`
	Purchase             *Purchase  `json:"-"`
}

type ListEntrySummary struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	Code             *string `json:"code" gorm:"unique"`
	ListName         *string `json:"listName"`
	AdditionalGuests uint    `json:"additionalGuests" gorm:"default:0"`
	ArrivalNote      *string `json:"arrivalNote"`
}

type ListEntrySummarySlice []ListEntrySummary

func (entry *ListEntry) MarkAsArrived() {
	now := time.Now()
	entry.ArrivedAt = &now
}

func (entries ListEntrySummarySlice) SortByQuery(q string) {
	query := strings.ToLower(q)
	sort.Slice(entries, func(i, j int) bool {
		nameI := strings.ToLower(entries[i].Name)
		nameJ := strings.ToLower(entries[j].Name)

		weightI := calculateWeight(nameI, query)
		weightJ := calculateWeight(nameJ, query)

		// Sort by weight first, then alphabetically
		if weightI != weightJ {
			return weightI > weightJ
		}
		return nameI < nameJ
	})
}

func calculateWeight(name, query string) int {
	switch {
	case name == query:
		return 3
	case strings.HasPrefix(name, query):
		return 2
	case strings.Contains(" "+name, " "+query):
		return 1
	default:
		return 0
	}
}

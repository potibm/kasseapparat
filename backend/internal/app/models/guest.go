package models

import (
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Guest represents a guest in a guestlist.
type Guest struct {
	GormOwnedModel

	GuestlistID          uint       `json:"guestlistId"`
	Guestlist            Guestlist  `json:"guestlist"`
	Name                 string     `json:"name"`
	Code                 *string    `gorm:"unique"               json:"code"`
	AdditionalGuests     uint       `gorm:"default:0"            json:"additionalGuests"`
	AttendedGuests       uint       `gorm:"default:0"            json:"attendedGuests"`
	ArrivedAt            *time.Time `json:"arrivedAt"`
	ArrivalNote          *string    `json:"arrivalNote"`
	NotifyOnArrivalEmail *string    `json:"notifyOnArrivalEmail"`
	PurchaseID           *uuid.UUID `json:"purchaseId"`
	Purchase             *Purchase  `json:"-"`
}

type GuestSummary struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	Code             *string `gorm:"unique"      json:"code"`
	ListName         *string `json:"listName"`
	AdditionalGuests uint    `gorm:"default:0"   json:"additionalGuests"`
	ArrivalNote      *string `json:"arrivalNote"`
}

type GuestSummarySlice []GuestSummary

func (entry *Guest) MarkAsArrived() {
	now := time.Now()
	entry.ArrivedAt = &now
}

func (entries GuestSummarySlice) SortByQuery(q string) {
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
	const (
		WheightExactMatch     = 3
		WheightPrefixMatch    = 2
		WheightSubstringMatch = 1
		WheightNoMatch        = 0
	)

	switch {
	case name == query:
		return WheightExactMatch
	case strings.HasPrefix(name, query):
		return WheightPrefixMatch
	case strings.Contains(" "+name, " "+query):
		return WheightSubstringMatch
	default:
		return WheightNoMatch
	}
}

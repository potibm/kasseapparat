package models

import (
	"sort"
	"strings"
)

// List represents a guestlist
type ListEntry struct {
	GormOwnedModel
	ListID 	uint	`json:"listId"`
	List 	List	`json:"list"`
	Name      string  `json:"name" `
	Code 		*string`json:"code" gorm:"unique"`
	AdditionalGuests uint `json:"additionalGuests" gorm:"default:0"`
	AttendedGuests uint `json:"attendedGuests" gorm:"default:0"`
	PurchaseID *uint `json:"purchaseId"`
	Purchase *Purchase `json:"-"`
}

type ListEntrySummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code 		*string`json:"code" gorm:"unique"`
	ListName *string `json:"listName"`
	AdditionalGuests uint `json:"additionalGuests" gorm:"default:0"`
}

type ListEntrySummarySlice []ListEntrySummary

func (entries ListEntrySummarySlice) SortByQuery(q string) {
	// Define custom sorting logic
	sort.Slice(entries, func(i, j int) bool {
		// Convert names and query to lowercase for case-insensitive comparison
		nameI, nameJ := strings.ToLower(entries[i].Name), strings.ToLower(entries[j].Name)
		query := strings.ToLower(q)

		// Prioritize exact matches on name
		if nameI == query && nameJ != query {
			return true
		}
		if nameJ == query && nameI != query {
			return false
		}

		// Then prioritize exact matches on code (not shown in summary)
		// You would need to modify the summary to include code for this part

		// Then prioritize names starting with the query
		if strings.HasPrefix(nameI, query) && !strings.HasPrefix(nameJ, query) {
			return true
		}
		if strings.HasPrefix(nameJ, query) && !strings.HasPrefix(nameI, query) {
			return false
		}

		// Then prioritize names containing the query as the start of any word
		if strings.Contains(" " + nameI, " " + query) && !strings.Contains(" " + nameJ, " " + query) {
			return true
		}
		if strings.Contains(" " + nameJ, " " + query) && !strings.Contains(" " + nameI, " " + query) {
			return false
		}

		// Finally, sort alphabetically
		return nameI < nameJ
	})
}
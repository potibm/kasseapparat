package storage

// GuestListFilters holds filtering criteria for guest lists, including a query string and a list of IDs.
type GuestListFilters struct {
	Query string
	IDs   []uint
}

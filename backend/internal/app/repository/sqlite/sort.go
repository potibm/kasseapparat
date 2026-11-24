package sqlite

import "errors"

// ErrInvalidSortField is returned by the list queries when the requested sort field is
// not one of the allowed columns. It is the client's mistake (react-admin sends _sort
// from the column headers), so handlers answer it with 400 instead of 500.
var ErrInvalidSortField = errors.New("invalid sort field name")

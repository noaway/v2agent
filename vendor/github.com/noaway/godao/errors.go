package godao

import "errors"

var (
	// ErrNoRowUpdated update affected no rows
	ErrNoRowUpdated = errors.New("no row updated")
)

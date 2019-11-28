package godao

import (
	"fmt"

	"strings"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// Transact execute function in the same transaction, rollback if error else commit
func Transact(db *gorm.DB, txFunc func(*gorm.DB) error) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
			tx.Rollback()
		}
	}()

	if err = txFunc(tx); err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}

// IsRecordNotFound check if select return nothing, only for postgres
func IsRecordNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}

// IsDuplicateEntry check if violate unique constraint, only for postgres
func IsDuplicateEntry(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		// pq unique_violation
		return err.Code == "23505"
	}
	return false
}

// IsTooManyConnections check if too many connections, only for postgres
func IsTooManyConnections(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		// pq too_many_connections
		return err.Code == "53300"
	}
	return false
}

// MustAffectedRows must update at least 1 row
func MustAffectedRows(res *gorm.DB) error {
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected <= 0 {
		return ErrNoRowUpdated
	}
	return nil
}

// CanRetry check if we can retry when error
func CanRetry(err error) bool {
	if err == nil {
		return false
	}

	if IsTooManyConnections(err) {
		return true
	}

	errMsg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(errMsg, "i/o timeout"):
		fallthrough
	case strings.Contains(errMsg, "no such host"):
		fallthrough
	case strings.Contains(errMsg, "connection refused"):
		return true
	default:
		return false
	}
}

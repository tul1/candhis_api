package model

import (
	"errors"
	"strings"
	"time"
)

type CandhisSessionID struct {
	id        string
	createdAt time.Time
}

func NewCandhisSessionID(id string, createdAt *time.Time) (CandhisSessionID, error) {
	if strings.HasPrefix(id, "PHPSESSID=") {
		return CandhisSessionID{}, errors.New("invalid session ID: contains PHPSESSID prefix")
	}
	if id == "" {
		return CandhisSessionID{}, errors.New("invalid session ID: cannot be empty")
	}

	if createdAt == nil {
		now := time.Now().UTC()
		createdAt = &now
	} else if createdAt.Location() != time.UTC {
		return CandhisSessionID{}, errors.New("invalid createdAt: must be in UTC format")
	}
	*createdAt = createdAt.Truncate(time.Microsecond) // I must truncate to avoid conflicts with database precision

	return CandhisSessionID{id: id, createdAt: *createdAt}, nil
}

func (c CandhisSessionID) ID() string {
	return c.id
}

func (c CandhisSessionID) PHPSESSID() string {
	return "PHPSESSID=" + c.id
}

func (c CandhisSessionID) CreatedAt() time.Time {
	return c.createdAt
}

// func (c CandhisSessionID) Value() (driver.Value, error) {
// 	// For the time.Time field, there’s no need to manually implement the Value() method because time.Time is natively supported by Go's SQL driver.
// 	// So, we don't need to add createdAt to this method.
// 	return c.id, nil
// }

package entities

import (
	"time"
)

// User is a user in the system.
type User struct {
	// ISU is the ITMO system user ID.
	ISU int64 `json:"isu"`
	// CalDavURL is the user's CalDav URL.
	CalDavURL string `json:"caldav_url"`
	// CreatedAt is the creation timestamp.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the last update timestamp.
	UpdatedAt time.Time `json:"updated_at"`
}

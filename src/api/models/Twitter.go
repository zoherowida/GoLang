package models

import (
	"errors"
	"html"
	"strings"
	"time"
)

type Twitter struct {
	ID        uint32    `gorm:"primary_id;auto_increment" json:"id"`
	Tweet     string    `gorm:"size:255;not null" json:"tweet"`
	CreatedAt time.Time `gorm:"default:current_timestamp()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
}

func (u *Twitter) Prepare() {
	u.ID = 0
	u.Tweet = html.EscapeString(strings.TrimSpace(u.Tweet))
}

// Validate validates the inputs
func (u *Twitter) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if u.Tweet == "" {
			return errors.New("Tweet is required")
		}
	case "update":
		if u.Tweet == "" {
			return errors.New("Tweet is required")
		}
	default:
		if u.Tweet == "" {
			return errors.New("Name is required")
		}
	}
	return nil
}

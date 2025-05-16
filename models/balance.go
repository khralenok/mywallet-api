package models

import "time"

type Balance struct {
	UserID       int       `json:"user_id"`
	Balance      int       `json:"balance"`
	SnapshotDate time.Time `json:"snapshot_date"`
}

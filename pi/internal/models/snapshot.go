package models

import "time"

type SnapShot struct {
	EventID    string
	FilePath   string
	CapturedAt time.Time
}

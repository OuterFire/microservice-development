package schema

import "time"

type NotificationMessage struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

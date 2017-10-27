package article

import (
	"time"
)

// Article describes the format of an article
type Article struct {
	ID        string `json:"_id,omitempty"`
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

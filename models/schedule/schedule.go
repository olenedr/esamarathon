package schedule

import "html/template"

// Schedule describes the structure of a parsed schedule
type Schedule struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Updated     string   `json:"updated,omitempty"`
	Link        string   `json:"link,omitempty"`
	Entries     []Entry  `json:"items,omitempty"`
	Columns     []string `json:"columns,omitempty"`
}

// Entry describes one speedrun in the schedule
type Entry struct {
	Scheduled string        `json:"scheduled,omitempty"`
	Game      template.HTML `json:"game,omitempty"`
	Estimate  string        `json:"estimate,omitempty"`
	Players   template.HTML `json:"players,omitempty"`
	Platform  string        `json:"platform,omitempty"`
	Category  string        `json:"category,omitempty"`
	Note      string        `json:"note,omitempty"`
	Data      []string      `json:"data,omitempty"`
	Length    float64       `json:"length_t,omitempty"`
}

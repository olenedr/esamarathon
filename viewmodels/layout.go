package viewmodels

import "github.com/esamarathon/website/models/menu"

type meta struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
}

// DefaultMata is a set of default metadata values
var DefaultMeta = meta{
	"ESA Marathon",
	"Welcome to European Speedrunner Assembly!",
	"http://www.esamarathon.com/static/img/og-image.png",
}

type layout struct {
	Meta meta      `json:"meta,omitempty"`
	Menu menu.Menu `json:"menu,omitempty"`
}

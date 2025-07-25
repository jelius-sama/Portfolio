package types

import "net/http"

// Type of Server Environment
type Environment struct {
	Prod string
	Dev  string
}

// Enum of Server ENvironment
var ENV = Environment{
	Prod: "production",
	Dev:  "development",
}

// Type representing reverse proxy status and settings
type BehindReverseProxy struct {
	StatementValid bool
	Port           string
}

// Type related to static page metadata
type MetaTag struct {
	Charset   string `json:"charset,omitempty"`
	Name      string `json:"name,omitempty"`
	Content   string `json:"content,omitempty"`
	Property  string `json:"property,omitempty"`
	HTTPEquiv string `json:"http-equiv,omitempty"`
}

type LinkTag struct {
	Rel         string `json:"rel"`
	Href        string `json:"href"`
	Type        string `json:"type,omitempty"`
	Crossorigin string `json:"crossorigin,omitempty"`
	Media       string `json:"media,omitempty"`
	Sizes       string `json:"sizes,omitempty"`
	As          string `json:"as,omitempty"`
	Referrer    string `json:"referrer,omitempty"`
	Title       string `json:"title,omitempty"`
}

type RouteMetadata struct {
	Path  string    `json:"path"`
	Title string    `json:"title,omitempty"`
	Meta  []MetaTag `json:"meta,omitempty"`
	Link  []LinkTag `json:"link,omitempty"`
}

type Middleware func(http.Handler) http.Handler

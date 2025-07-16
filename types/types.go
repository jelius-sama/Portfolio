package types

type Environment struct {
	Prod string
	Dev  string
}

var ENV = Environment{
	Prod: "production",
	Dev:  "development",
}

type MetaTag struct {
	Charset   string `json:"charset,omitempty"`
	Name      string `json:"name,omitempty"`
	Content   string `json:"content,omitempty"`
	Property  string `json:"property,omitempty"`
	HTTPEquiv string `json:"http-equiv,omitempty"`
}

type LinkTag struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type RouteMetadata struct {
	Path  string    `json:"path"`
	Title string    `json:"title,omitempty"`
	Meta  []MetaTag `json:"meta,omitempty"`
	Link  []LinkTag `json:"link,omitempty"`
}

package types

type MLink struct {
    Rel   string  `json:"rel"`
    Href  string  `json:"href"`
    Media *string `json:"media"`
}

type MMeta struct {
    Name     *string `json:"name"`
    Property *string `json:"property"`
    Content  string  `json:"content"`
}

type Metadata struct {
    Path        string  `json:"path"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Links       []MLink `json:"links"`
    Meta        []MMeta `json:"meta"`
}


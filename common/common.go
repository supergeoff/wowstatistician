package common

// Links struct format
type Links struct {
	Self URL `json:"self"`
}

// URL struct format
type URL struct {
	Href string `json:"href"`
}

// Value struct format
type Value struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Media struct format
type Media struct {
	Key URL `json:"key"`
	ID  int `json:"id"`
}

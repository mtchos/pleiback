package service

type Track struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Href    string    `json:"href"`
	URI     string    `json:"uri"`
	Artists []*Artist `json:"artists"`
}

type TrackFilter struct {
	Query  string `json:"query"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

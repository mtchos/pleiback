package service

type Artist struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Href   string    `json:"href"`
	URI    string    `json:"uri"`
	Genres []*string `json:"genres"`
}

package model

type Module struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Sections []Section `json:"sections"`
}

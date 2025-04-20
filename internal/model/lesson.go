package model

type Lesson struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	XP    int64  `json:"xp"`
}

package model

type Section struct {
	ID      int64    `json:"id"`
	Title   string   `json:"title"`
	Lessons []Lesson `json:"lessons"`
}

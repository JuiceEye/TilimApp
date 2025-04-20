package model

type Section struct {
	ID      int64
	Title   string
	Lessons []Lesson
}

package model

type Answer struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Image     string `json:"image"`
	IsCorrect string `json:"is_correct"`
}

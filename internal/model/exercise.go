package model

type Exercise struct {
	ID           int64    `json:"id"`
	Text         string   `json:"text"`
	Image        string   `json:"image"`
	QuestionText string   `json:"question_text"`
	Answers      []Answer `json:"answers"`
}

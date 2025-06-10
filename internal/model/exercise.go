package model

type Exercise struct {
	ID           int64    `json:"id"`
	TypeCode     string   `json:"type_code"`
	Text         string   `json:"text"`
	Image        string   `json:"image"`
	QuestionText string   `json:"question_text"`
	Audio        *File    `json:"audio"`
	Answers      []Answer `json:"answers"`
}

type File struct {
	UUID string `json:"uuid"`
	Body string `json:"body"`
}

package response

type GetMainPageModuleResponse struct {
	ID       int64             `json:"user_id"`
	Title    string            `json:"title"`
	Sections []MainPageSection `json:"sections"`
}

type MainPageSection struct {
	ID      string           `json:"id"`
	Title   string           `json:"title"`
	Lessons []MainPageLesson `json:"lessons"`
}

type MainPageLesson struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	XP    int64  `json:"xp"`
}

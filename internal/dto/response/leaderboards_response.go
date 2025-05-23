package response

type GetLeaderboardsResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	XPPoints int64  `json:"xp_points"`
	Image    string `json:"image"`
}

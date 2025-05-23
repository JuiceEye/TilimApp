package response

import "tilimauth/internal/model"

type GetLeaderboardsResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	XPPoints int64  `json:"xp_points"`
	Image    string `json:"image"`
}

func ToLeaderboardsResponseList(profiles []*model.LeaderboardsUser) []GetLeaderboardsResponse {
	leaderboards := make([]GetLeaderboardsResponse, 0, len(profiles))
	for _, p := range profiles {
		leaderboards = append(leaderboards, GetLeaderboardsResponse{
			UserID:   p.UserID,
			Username: p.Username,
			XPPoints: p.XPPoints,
			Image:    p.Image,
		})
	}
	return leaderboards
}

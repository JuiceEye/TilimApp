package response

import "tilimauth/internal/model"

type GetLeaderboardsProfileResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	XPPoints int64  `json:"xp_points"`
	Image    string `json:"image"`
}

func ToProfileResponseList(profiles []*model.Profile) []GetLeaderboardsProfileResponse {
	leaderboards := make([]GetLeaderboardsProfileResponse, 0, len(profiles))
	for _, p := range profiles {
		leaderboards = append(leaderboards, GetLeaderboardsProfileResponse{
			UserID:   p.UserID,
			Username: p.Username,
			XPPoints: p.XPPoints,
			Image:    p.Image,
		})
	}
	return leaderboards
}

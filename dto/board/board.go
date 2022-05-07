package boardDto

type SetClientScoreInput struct {
	Score float64 `json:"score" binding:"min=0,max=100"`
}

type ClientAndScoreSet struct {
	ClientId string  `json:"clientId"`
	Score    float64 `json:"score"`
}

type GetClientLeaderboaredInput struct {
	Top int64 `form:"top"`
}

type GetClientLeaderboardOutput struct {
	TopPlayers []ClientAndScoreSet `json:"topPlayers"`
}

package boardService

import boardDomain "leaderboard/domain/board"

type IBoardService interface {
	SetScore(clientId string, score float64) error
	GetLeaderboard(top int64) ([]boardDomain.MemberScore, error)
}

type MemberScore struct {
	Member string
	Score  float64
}

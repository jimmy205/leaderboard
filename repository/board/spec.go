package boardRepository

import (
	"time"

	boardDomain "leaderboard/domain/board"
)

type IBoardRepository interface {
	BoardKeyExist() (bool, error)
	SetBoardExpire(expire time.Duration) error
	SetScore(clientId string, score float64) error
	GetLeaderboard(top int64) ([]boardDomain.MemberScore, error)
}

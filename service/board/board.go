package boardService

import (
	"time"

	boardDomain "leaderboard/domain/board"
	boardRepository "leaderboard/repository/board"
)

type boardService struct {
	boardRepo boardRepository.IBoardRepository
}

func NewBoardService(
	boardRepo boardRepository.IBoardRepository,
) IBoardService {
	return boardService{
		boardRepo: boardRepo,
	}
}

const (
	expireTime = time.Minute * 10
)

// will set expire time to 10mins if board is not exist
func (s boardService) SetScore(clientId string, score float64) error {
	exist, err := s.boardRepo.BoardKeyExist()
	if err != nil {
		return err
	}
	if err := s.boardRepo.SetScore(clientId, score); err != nil {
		return err
	}
	if !exist {
		if err := s.boardRepo.SetBoardExpire(expireTime); err != nil {
			return err
		}
	}
	return nil
}

func (s boardService) GetLeaderboard(top int64) ([]boardDomain.MemberScore, error) {
	list, err := s.boardRepo.GetLeaderboard(top)
	if err != nil {
		return nil, err
	}
	return list, nil
}

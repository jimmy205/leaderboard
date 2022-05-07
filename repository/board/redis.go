package boardRepository

import (
	"context"
	"errors"
	"time"

	boardDomain "leaderboard/domain/board"

	"github.com/go-redis/redis/v8"
)

const (
	boardKey = "leaderboard"
)

type redisBoard struct {
	cmd *redis.Client
}

func NewBoardRedisRepositroy(cmd *redis.Client) IBoardRepository {
	return redisBoard{
		cmd: cmd,
	}
}

func (r redisBoard) SetScore(clientId string, score float64) error {
	cmd := r.cmd.ZAdd(context.Background(), boardKey, &redis.Z{
		Member: clientId,
		Score:  score,
	})
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r redisBoard) GetLeaderboard(top int64) ([]boardDomain.MemberScore, error) {
	cmd := r.cmd.ZRevRangeWithScores(context.Background(), boardKey, 0, top-1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	list := []boardDomain.MemberScore{}
	for _, v := range cmd.Val() {
		name, ok := v.Member.(string)
		if !ok {
			return nil, errors.New("redis assert username failed.")
		}
		list = append(list, boardDomain.MemberScore{
			Member: name,
			Score:  v.Score,
		})
	}
	return list, nil
}

func (r redisBoard) SetBoardExpire(expire time.Duration) error {
	cmd := r.cmd.Expire(context.Background(), boardKey, expire)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r redisBoard) BoardKeyExist() (bool, error) {
	cmd := r.cmd.Exists(context.Background(), boardKey)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	// value = 1 is true
	exist := false
	if cmd.Val() == 1 {
		exist = true
	}
	return exist, nil
}

package boardUsecase

import (
	boardDto "leaderboard/dto/board"
	"leaderboard/dto/response"
	boardService "leaderboard/service/board"

	"github.com/gin-gonic/gin"
)

type IBoardUsecase interface {
	SetScore(ctx *gin.Context)
	GetLeaderboard(ctx *gin.Context)
}

type boardUsecase struct {
	boardService boardService.IBoardService
}

func NewBoardController(
	service boardService.IBoardService,
) IBoardUsecase {
	return boardUsecase{
		boardService: service,
	}
}

const (
	headerClientId = "clientId"
)

func (u boardUsecase) SetScore(ctx *gin.Context) {
	input := boardDto.SetClientScoreInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.FailWithParam(ctx)
		return
	}

	clientId := ctx.GetHeader(headerClientId)
	if clientId == "" {
		response.FailWithParam(ctx)
		return
	}

	if err := u.boardService.SetScore(clientId, input.Score); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.SuccessWithStatusOk(ctx)
}

const (
	top10 = 10
)

func (u boardUsecase) GetLeaderboard(ctx *gin.Context) {
	list, err := u.boardService.GetLeaderboard(top10)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	players := []boardDto.ClientAndScoreSet{}
	for _, l := range list {
		players = append(players, boardDto.ClientAndScoreSet{
			ClientId: l.Member,
			Score:    l.Score,
		})
	}

	response.SuccessWithData(ctx, boardDto.GetClientLeaderboardOutput{
		TopPlayers: players,
	})
	return
}

package router

import (
	boardUsecase "leaderboard/usecase/board"

	"github.com/gin-gonic/gin"
)

func setBoard(server *gin.Engine, usecase boardUsecase.IBoardUsecase) {
	server.GET("/api/v1/leaderboard", usecase.GetLeaderboard)
	server.POST("/api/v1/score", usecase.SetScore)
}

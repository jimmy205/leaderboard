package router

import (
	"os"

	boardUsecase "leaderboard/usecase/board"

	"github.com/gin-gonic/gin"
)

func StartRouter(
	boardUsecase boardUsecase.IBoardUsecase,
) *gin.Engine {
	server := gin.Default()

	setBoard(server, boardUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	return server
}

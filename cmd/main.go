package main

import (
	"os"

	"leaderboard/lib"
	boardRepository "leaderboard/repository/board"
	"leaderboard/router"
	boardService "leaderboard/service/board"
	boardUsecase "leaderboard/usecase/board"
)

func main() {
	boardRepository := boardRepository.NewBoardRedisRepositroy(lib.GetRedis())
	boardService := boardService.NewBoardService(boardRepository)
	boardUsecase := boardUsecase.NewBoardController(boardService)

	server := router.StartRouter(boardUsecase)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	server.Run(":" + port)
}

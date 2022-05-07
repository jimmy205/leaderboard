package boardUsecase_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	boardDto "leaderboard/dto/board"
	boardRepository "leaderboard/repository/board"
	"leaderboard/router"
	boardService "leaderboard/service/board"
	boardUsecase "leaderboard/usecase/board"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type board struct {
	server  *gin.Engine
	mockDb  *redis.Client
	usecase boardUsecase.IBoardUsecase
}

var boardTest board

func TestMain(m *testing.M) {
	mockDb, err := miniredis.Run()
	if err != nil {
		return
	}

	mock := redis.NewClient(&redis.Options{
		Addr: mockDb.Addr(),
	})

	usecase := boardUsecase.NewBoardController(
		boardService.NewBoardService(
			boardRepository.NewBoardRedisRepositroy(mock),
		),
	)

	boardTest = board{
		mockDb:  mock,
		server:  router.StartRouter(usecase),
		usecase: usecase,
	}

	m.Run()
}

type arg struct {
	method       string
	path         string
	header       map[string]string
	params       interface{}
	expectedCode int
	expectedBody string
}

const (
	methodPOST = "POST"
	methodGET  = "GET"

	pathSetScore = "/api/v1/score"
	pathGetBoard = "/api/v1/leaderboard"

	responseBodyOK      = `{"status":"OK"}`
	responseParamsError = `{"error":"invalid params"}`
)

func (arg *arg) runTest() error {
	writer := httptest.NewRecorder()
	b, err := json.Marshal(arg.params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(arg.method, arg.path, bytes.NewReader(b))
	if err != nil {
		return err
	}

	for k, v := range arg.header {
		req.Header.Set(k, v)
	}

	boardTest.server.ServeHTTP(writer, req)

	if arg.expectedCode != writer.Code {
		return fmt.Errorf("response code error, want: %d, got: %d", arg.expectedCode, writer.Code)
	}
	if arg.expectedBody != writer.Body.String() {
		return fmt.Errorf("response body error, want: %s, got: %s", arg.expectedBody, writer.Body.String())
	}

	return nil
}

func (arg arg) formatError(err error) string {
	return fmt.Sprintf("*** Testing Error *** :" + err.Error())
}

func Test_SetScoreNormal(t *testing.T) {
	arg := &arg{
		method: methodPOST,
		path:   pathSetScore,
		header: map[string]string{"ClientId": "Jimmy"},
		params: boardDto.SetClientScoreInput{
			Score: 11.2,
		},
		expectedCode: http.StatusOK,
		expectedBody: responseBodyOK,
	}

	if err := arg.runTest(); err != nil {
		t.Error(arg.formatError(err))
	}
}

func Test_SetScoreMissingHeader(t *testing.T) {
	arg := &arg{
		method: methodPOST,
		path:   pathSetScore,
		params: boardDto.SetClientScoreInput{
			Score: 11.2,
		},
		expectedCode: http.StatusBadRequest,
		expectedBody: responseParamsError,
	}

	if err := arg.runTest(); err != nil {
		t.Error(arg.formatError(err))
	}
}

func Test_SetScoreNegative(t *testing.T) {
	arg := &arg{
		method: methodPOST,
		path:   pathSetScore,
		header: map[string]string{"ClientId": "Jimmy"},
		params: boardDto.SetClientScoreInput{
			Score: -1,
		},
		expectedCode: http.StatusBadRequest,
		expectedBody: responseParamsError,
	}

	if err := arg.runTest(); err != nil {
		t.Error(arg.formatError(err))
	}
}

func Test_SetScoreOver100(t *testing.T) {
	arg := &arg{
		method: methodPOST,
		path:   pathSetScore,
		header: map[string]string{"ClientId": "Jimmy"},
		params: boardDto.SetClientScoreInput{
			Score: 1002,
		},
		expectedCode: http.StatusBadRequest,
		expectedBody: responseParamsError,
	}

	if err := arg.runTest(); err != nil {
		t.Error(arg.formatError(err))
	}
}

func Test_NormalGetBoard(t *testing.T) {
	// 先把之前的測資都刪掉
	boardTest.mockDb.Del(context.Background(), "leaderboard")

	palyersMap := map[string]float64{
		"player1": 100,
		"player2": 95,
		"player3": 44,
		"player4": 22,
		"player5": 33,
	}

	playerReq := []*redis.Z{}
	playerRes := []boardDto.ClientAndScoreSet{}
	for k, v := range palyersMap {
		playerReq = append(playerReq, &redis.Z{Member: k, Score: v})
		playerRes = append(playerRes, boardDto.ClientAndScoreSet{ClientId: k, Score: v})
	}

	boardTest.mockDb.ZAdd(context.Background(), "leaderboard", playerReq...)

	sort.Slice(playerRes, func(i, j int) bool { return playerRes[i].Score > playerRes[j].Score })

	b, _ := json.Marshal(boardDto.GetClientLeaderboardOutput{TopPlayers: playerRes})
	arg := &arg{
		method: methodGET,
		path:   pathGetBoard,
		params: boardDto.GetClientLeaderboaredInput{
			Top: 5,
		},
		expectedCode: http.StatusOK,
		expectedBody: string(b),
	}

	if err := arg.runTest(); err != nil {
		t.Error(arg.formatError(err))
	}
}

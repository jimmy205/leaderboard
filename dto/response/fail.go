package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrParams = "invalid params"
)

type FailResponse struct {
	Error string `json:"error,omitempty"`
}

func FailWithParam(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, FailResponse{
		Error: ErrParams,
	})
}

func FailWithError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, FailResponse{
		Error: err.Error(),
	})
}

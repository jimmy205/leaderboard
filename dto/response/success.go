package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessWithStatus struct {
	Status string `json:"status"`
}

func SuccessWithData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func SuccessWithStatusOk(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, SuccessWithStatus{
		Status: "OK",
	})
}

package apis

import (
	"net/http"

	"github.com/cylonchau/firewalldGateway/server/apis"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `form:"apis" json:"apis,omitempty"`
	Msg  string      `form:"msg" json:"msg,omitempty"`
	Data interface{} `form:"data" json:"data,omitempty"`
}

// APIResponse ....
func APIResponse(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := apis.DecodeErr(err)
	ctx.JSON(http.StatusOK, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// SuccessResponse ....
func SuccessResponse(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := apis.DecodeErr(err)
	ctx.JSON(http.StatusOK, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// NotFountResponse ....
func NotFount(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := apis.DecodeErr(err)
	ctx.JSON(http.StatusNotFound, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// ConnectDbusService ....
func ConnectDbusService(ctx *gin.Context, err error) {
	returnCode, message := apis.DecodeErr(err)
	ctx.JSON(http.StatusInternalServerError, Response{
		Code: returnCode,
		Msg:  message,
		Data: apis.ErrDBus,
	})
}

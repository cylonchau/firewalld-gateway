package query

import (
	"firewall-api/code"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `form:"code" json:"code,omitempty"`
	Msg  string      `form:"msg" json:"msg,omitempty"`
	Data interface{} `form:"data" json:"data,omitempty"`
}

// APIResponse ....
func APIResponse(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := code.DecodeErr(err)
	ctx.JSON(http.StatusOK, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// SuccessResponse ....
func SuccessResponse(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := code.DecodeErr(err)
	ctx.JSON(http.StatusOK, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// NotFountResponse ....
func NotFount(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := code.DecodeErr(err)
	ctx.JSON(http.StatusNotFound, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// ConnectDbusService ....
func ConnectDbusService(ctx *gin.Context, err error) {
	returnCode, message := code.DecodeErr(err)
	ctx.JSON(http.StatusInternalServerError, Response{
		Code: returnCode,
		Msg:  message,
		Data: code.ErrDBus,
	})
}

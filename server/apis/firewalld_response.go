package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `form:"code" json:"code,omitempty"`
	Msg  string      `form:"msg" json:"msg,omitempty"`
	Data interface{} `form:"data" json:"data,omitempty"`
}

type ResponseSlice struct {
	Code int         `form:"code" json:"code,omitempty"`
	Msg  []string    `form:"msg" json:"msg,omitempty"`
	Data interface{} `form:"data" json:"data,omitempty"`
}

func AuthFailed(ctx *gin.Context, msg *Errno, data interface{}) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Code: msg.Code,
		Msg:  msg.Message,
		Data: data,
	})
}

// API403Response ....
func Auth403Failed(ctx *gin.Context, msg *Errno, data interface{}) {
	ctx.JSON(http.StatusForbidden, Response{
		Code: msg.Code,
		Msg:  msg.Message,
		Data: data,
	})
}

// APIResponse ....
func APIResponse(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(http.StatusOK, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// 404Response
func API404Response(ctx *gin.Context, err error) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(http.StatusNotFound, Response{
		Code: returnCode,
		Msg:  message,
	})
}

// 409Response
func API409Response(ctx *gin.Context, err error) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(http.StatusConflict, Response{
		Code: returnCode,
		Msg:  message,
	})
}

// 400Response
func API500Response(ctx *gin.Context, err error) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(http.StatusInternalServerError, Response{
		Code: returnCode,
		Msg:  message,
	})
}

// SuccessResponse ....
func SuccessResponse(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(http.StatusOK, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// NotFountResponse ....
func NotFount(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(http.StatusNotFound, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// ConnectDbusService ....
func ConnectDbusService(ctx *gin.Context, err error) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(http.StatusInternalServerError, Response{
		Code: returnCode,
		Msg:  message,
		Data: ErrDBus,
	})
}

// ConnectDbusService ....
func BacthMissionFailedResponse(ctx *gin.Context, err ...error) {
	returnCode, message := DecodeErrSlice(err...)
	ctx.JSON(http.StatusAccepted, ResponseSlice{
		Code: returnCode,
		Msg:  message,
		Data: BatchErrCreated,
	})
}

// ConnectDbusService ....
func BacthMissionSuccessResponse(ctx *gin.Context, err error) {
	returnCode, _ := DecodeErr(err)
	ctx.JSON(http.StatusCreated, ResponseSlice{
		Code: returnCode,
		Data: BatchSuccessCreated,
	})
}

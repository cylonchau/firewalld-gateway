package apis

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

//nolint: golint
var (
	// Common errors
	OK                         = &Errno{Code: 10000, Message: "operation succeeded"}
	NETWORK_MASQUERADE_ENABLE  = &Errno{Code: 0, Message: "network masquerade is enable"}
	NETWORK_MASQUERADE_DISABLE = &Errno{Code: 0, Message: "network masquerade is disable"}
	ErrDBus                    = "connect to remote firewalld server failed"
	InternalServerError        = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind                    = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct"}
	ErrParam                   = &Errno{Code: 10003, Message: "参数有误"}
	ErrSignParam               = &Errno{Code: 10004, Message: "签名参数有误"}

	ErrValidation         = &Errno{Code: 20001, Message: "Validation failed"}
	ErrDatabase           = &Errno{Code: 20002, Message: "Database error"}
	ErrToken              = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token"}
	ErrInvalidTransaction = &Errno{Code: 20004, Message: "invalid transaction"}

	// NOTFOUNT
	ErrRichNotFount    = &Errno{Code: 40004, Message: "The rich rules in the zone is empty"}
	ErrServiceNotFount = &Errno{Code: 40004, Message: "The service in the zone is empty"}
	ErrPortNotFount    = &Errno{Code: 40004, Message: "The port in the zone is empty"}
	ErrZoneNotFount    = &Errno{Code: 40004, Message: "Not found the zone"}
	ErrForwardNotFount = &Errno{Code: 40004, Message: "The Forward in the zone is empty"}

	// auther errors
	ErrEncrypt               = &Errno{Code: 50101, Message: "success"}
	ErrUserNotFound          = &Errno{Code: 50102, Message: "User not found"}
	ErrTokenInvalid          = &Errno{Code: 50103, Message: "Invalied token"}
	ErrPasswordIncorrect     = &Errno{Code: 50104, Message: "Incorrect username or password"}
	ErrUserExist             = &Errno{Code: 50105, Message: "User exists"}
	ErrUserNotExist          = &Errno{Code: 50106, Message: "User does not exist"}
	ErrNeedAuth              = &Errno{Code: 50107, Message: "Your need authetication"}
	ErrSendSMSTooMany        = &Errno{Code: 50109, Message: "已超出当日限制，请明天再试"}
	ErrVerifyCode            = &Errno{Code: 50110, Message: "验证码错误"}
	ErrEmailOrPassword       = &Errno{Code: 50111, Message: "邮箱或密码错误"}
	ErrTwicePasswordNotMatch = &Errno{Code: 50112, Message: "两次密码输入不一致"}
	ErrRegisterFailed        = &Errno{Code: 50113, Message: "注册失败"}
	ErrCreatedUser           = &Errno{Code: 50114, Message: "用户创建失败"}

	ErrTagNotFound  = &Errno{Code: 30104, Message: "Tag not found"}
	ErrUserTagExist = &Errno{Code: 30106, Message: "Tag does existed"}

	BatchSuccessCreated = &Errno{Code: 60000, Message: "The batch mission has created"}
	BatchErrCreated     = &Errno{Code: 60005, Message: "The batch mission create failed"}
)

// Errno ...
type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

// Err represents an error
type Err struct {
	Code    int
	Message string
	Err     error
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - apis: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

// DecodeErr ...
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	case error:
		err = errors.New(strings.ReplaceAll(err.Error(), "\"", ""))
	default:
	}

	return InternalServerError.Code, err.Error()
}

func DecodeErrSlice(err ...error) (int, []string) {
	var returnStr []string
	for _, n := range err {
		if n == nil {
			returnStr = append(returnStr, OK.Message)
			continue
		}
		switch typed := n.(type) {
		case *Err:
			returnStr = append(returnStr, typed.Message)
			continue
		case *Errno:
			returnStr = append(returnStr, typed.Message)
			continue
		default:
		}
		returnStr = append(returnStr, n.Error())
	}
	return http.StatusAccepted, returnStr
}

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1. 定义最基础的 Code
type resCode int

const (
	CodeSuccess      resCode = 10000
	CodeInvalidParam resCode = 10001
	CodeServerBusy   resCode = 10002
	CodeNeedLogin    resCode = 10003
)

// 2. 简化的消息映射
var codeMsgMap = map[resCode]string{
	CodeSuccess:      "success",
	CodeInvalidParam: "请求参数错误",
	CodeServerBusy:   "服务器繁忙",
	CodeNeedLogin:    "请先登录",
}

type ResponseData struct {
	Code resCode `json:"code"`
	Msg  string  `json:"msg"` // 通常 Msg 是字符串，不用 any
	Data any     `json:"data"`
}

// 3. 核心：统一响应入口 (私有方法，减少重复代码)
func response(c *gin.Context, code resCode, data any) {
	msg, ok := codeMsgMap[code]
	if !ok {
		msg = codeMsgMap[CodeServerBusy] // 默认兜底错误
	}

	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success 成功 (外部直接调用这个)
func Success(c *gin.Context, data any) {
	response(c, CodeSuccess, data)
}

// Fail 失败 (外部直接调用这个)
func Fail(c *gin.Context, code resCode) {
	// 失败时 Data 给 nil，比给 "" 更安全，兼容性更好
	response(c, code, nil)
}

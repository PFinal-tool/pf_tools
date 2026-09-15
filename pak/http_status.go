package pak

import (
	"fmt"
	"strconv"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: HTTP 状态码速查
 * @Project: pf_tools
 */

var httpStatusText = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	200: "OK",
	201: "Created",
	202: "Accepted",
	204: "No Content",
	301: "Moved Permanently",
	302: "Found",
	304: "Not Modified",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	413: "Payload Too Large",
	415: "Unsupported Media Type",
	418: "I'm a teapot",
	422: "Unprocessable Entity",
	429: "Too Many Requests",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
}

var httpStatusHint = map[int]string{
	200: "请求成功",
	201: "资源已创建",
	204: "成功但无响应体",
	301: "永久重定向",
	302: "临时重定向",
	304: "缓存未变化",
	400: "请求参数/格式错误",
	401: "未认证（缺/错 Token）",
	403: "已认证但无权限",
	404: "资源不存在",
	409: "资源冲突",
	422: "语义错误（校验失败）",
	429: "限流/请求过频",
	500: "服务端内部错误",
	502: "网关上游异常",
	503: "服务不可用（过载/维护）",
	504: "网关等待上游超时",
}

type StatusTool struct {
	Quiet bool
}

func (s *StatusTool) Run(codeStr string) {
	view := ResultView{Quiet: s.Quiet}
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		view.ShowText("HTTP状态码", "请输入数字状态码，如 404", "", false)
		return
	}
	textEN, ok := httpStatusText[code]
	if !ok {
		textEN = "未知/未收录"
	}
	hint := httpStatusHint[code]
	if hint == "" {
		hint = categoryHint(code)
	}
	msg := fmt.Sprintf("%d %s\n说明: %s", code, textEN, hint)
	view.ShowText("HTTP状态码", msg, fmt.Sprintf("%d %s", code, textEN), ok)
}

func categoryHint(code int) string {
	switch {
	case code >= 100 && code < 200:
		return "信息响应"
	case code >= 200 && code < 300:
		return "成功响应"
	case code >= 300 && code < 400:
		return "重定向"
	case code >= 400 && code < 500:
		return "客户端错误"
	case code >= 500 && code < 600:
		return "服务端错误"
	default:
		return "非标准状态码"
	}
}

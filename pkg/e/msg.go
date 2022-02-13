package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",
	ERROR_AUTH:     "认证失败",

	ERROR_CMD_REGEXP:         "命令含有特殊字符",
	ERROR_CMD_USER_NOT_FOUND: "runas用户不存在",
	ERROR_PARSE_CMD_ARGS:     "命令参数不合法",
	ERROR_SET_CMD_RUNAS:      "设置用户执行时失败",
	ERROR_CMD_RUN:            "命令执行失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}

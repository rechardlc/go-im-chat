package xeer

var CodeText = map[int]string{
	SERVER_COMMON_ERROR: "服务器异常，请稍后再试",
	REQUEST_PARAM_ERROR: "请求参数错误，请检查参数",
	DB_ERROR: "数据库繁忙，请稍后再试",
}

func Errmsg(errcode int) string {
	if text, ok := CodeText[errcode]; ok {
		return text
	}
	return CodeText[SERVER_COMMON_ERROR]
}
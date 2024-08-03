package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR:  "Server Error, Try Later! ",
	REQUEST_PARAM_ERRROR: "Req Param Error! ",
	DB_ERROR:             "Database Busy, Try Later! ",
}

func ErrMsg(errCode int) string {
	if msg, ok := codeText[errCode]; ok {
		return msg
	}
	return codeText[SERVER_COMMON_ERROR]
}

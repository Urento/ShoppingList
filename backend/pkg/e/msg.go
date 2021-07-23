package e

var MsgFlags = map[int]string{
	SUCCESS:                   "ok",
	ERROR:                     "fail",
	INVALID_PARAMS:            "request parameter error",
	ERROR_NOT_AUTHORIZED:      "not authorized to access this route",
	ERROR_RATELIMIT_TRY_LATER: "ratelimit reached; try again in 1 minute",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

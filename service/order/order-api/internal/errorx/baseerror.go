package errorx

const (
	defaultErrCode = 1001
	RPCErrCode     = 1002
	MySQLErrCode   = 1003
	ErrCode        = 1004
)

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CodeErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) Error() string {
	return e.Msg
}

func NewCodeError(code int, msg string) error {
	return &CodeError{
		Code: code,
		Msg:  msg,
	}
}

func NewDefaultCodeError(msg string) error {
	return &CodeError{
		Code: defaultErrCode,
		Msg:  msg,
	}
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

package middleware

import (
	"bytes"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// 定义全局中间件

// 功能：
// 记录所有请求的响应信息

type bodyCopy struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func NewBodyCopy(w http.ResponseWriter) *bodyCopy {
	return &bodyCopy{
		ResponseWriter: w,
		body:           bytes.NewBuffer([]byte{}),
	}
}

func (bc bodyCopy) Write(b []byte) (int, error) {
	// 1. 先在我的小本本记录响应内容
	bc.body.Write(b)
	// 2. 再往HTTP响应里写响应内容
	return bc.ResponseWriter.Write(b)
}

func CopyResp(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bc := NewBodyCopy(w)
		next(bc, r)
		logx.Debugf("--> req: %v resp: %v\n", r.URL, bc.body.String())
	}
}

func MiddlewareWithAnotherService(ok bool) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if ok {
				logx.Debug("Another service is enabled")
			}
			next(w, r)
		}
	}
}

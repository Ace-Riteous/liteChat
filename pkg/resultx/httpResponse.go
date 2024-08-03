package resultx

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"liteChat/pkg/xerr"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: http.StatusOK,
		Msg:  "",
		Data: data,
	}
}

func Fail(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}

func ErrHandler(serviceName string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errCode := xerr.SERVER_COMMON_ERROR
		errMsg := xerr.ErrMsg(errCode)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zerr.CodeMsg); ok {
			errCode = e.Code
			errMsg = e.Msg
		} else {
			if gstatus, ok := status.FromError(err); ok {
				errCode = int(gstatus.Code())
				errMsg = gstatus.Message()
			}
		}

		logx.WithContext(ctx).Errorf("【%s】 err %v", serviceName, err)

		return http.StatusBadRequest, Fail(errCode, errMsg)
	}
}

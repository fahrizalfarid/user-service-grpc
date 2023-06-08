package wrapper

import (
	"context"
	"sync"

	"github.com/fahrizalfarid/user-service-grpc/src/constant"
	"github.com/fahrizalfarid/user-service-grpc/utils"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
)

type Wrapper struct {
	Authorization utils.Authentication
	Waitgroup     *sync.WaitGroup
}

func NewWrapper(auth utils.Authentication, wg *sync.WaitGroup) *Wrapper {
	return &Wrapper{
		Authorization: auth,
		Waitgroup:     wg,
	}
}

func (w *Wrapper) Authentication(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		md, ok := metadata.FromContext(ctx)
		if !ok {
			return constant.ErrAuth
		}

		token, exist := md.Get("token")
		if !exist {
			return constant.ErrAuth
		}

		_, err := w.Authorization.ParsingToken(token)
		if err != nil {
			return err
		}

		err = fn(ctx, req, rsp)
		return err
	}
}

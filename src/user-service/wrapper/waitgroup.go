package wrapper

import (
	"context"

	"go-micro.dev/v4/server"
)

func (w *Wrapper) WaitGroup() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			w.Waitgroup.Add(1)
			defer w.Waitgroup.Done()
			return fn(ctx, req, rsp)
		}
	}
}

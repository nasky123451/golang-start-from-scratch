package tracing

import (
	"context"
	"time"
)

func doOperationWithCtx(ctx context.Context) error {
	// 模擬操作
	time.Sleep(100 * time.Millisecond)
	return nil
}

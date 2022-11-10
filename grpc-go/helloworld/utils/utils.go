package utils

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
)

// ConditionOperation : retry decision depends on the returned bool
type ConditionOperation func() (bool, error)


func RetryWithCondition(ctx context.Context, b backoff.BackOff, o ConditionOperation) error {
	ticker := backoff.NewTicker(b)
	defer ticker.Stop()
	var err error
	var needRetry bool
	for {
		select {
		case <-ctx.Done():
			return errors.Wrapf(ctx.Err(), "stopped retrying err: %v", err)
		default:
			select {
			case _, ok := <-ticker.C:
				if !ok {
					return err
				}
				needRetry, err = o()
				if !needRetry {
					return err
				}
			case <-ctx.Done():
				return errors.Wrapf(ctx.Err(), "stopped retrying err: %v", err)
			}
		}
	}
}

func HandleRequest(ctx context.Context) error {
	doneChan := make(chan struct{}, 1)
	errChan := make(chan error, 1)
	// 此处可以填充处理逻辑，doneChan 和 errChan 可以接收处理结果

	select {
	case <- doneChan:
		return nil
	case err := <- errChan:
		return err
	case <- ctx.Done():
		return ctx.Err()
	}
}


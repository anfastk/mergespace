package kafka

import (
	"context"
	"time"
)

func retry(ctx context.Context,max int,fn func() error) error {

	var err error

	for i := 0; i < max; i++ {

		select {

		case <-ctx.Done():
			return ctx.Err()

		default:
		}

		err = fn()

		if err == nil {
			return nil
		}

		time.Sleep(
			time.Duration(i+1) * time.Second,
		)
	}

	return err
}
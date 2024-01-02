package batch

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/goark/errs"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

var ErrInterrupt = errs.New("interrupted")

type Batch interface {
	Run(ctx context.Context) error
	Name() string
}

type Done struct{}

func RunBatch(batch Batch) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(cancelCtx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan Done, 1)

	eg.Go(func() error {
		if err := batch.Run(ctx); err != nil {
			return errs.Wrap(err)
		}
		done <- Done{}

		return nil
	})

	eg.Go(func() error {
		select {
		case <-quit:
			defer cancel()
			return ErrInterrupt
		case <-ctx.Done():
			return nil
		case <-done:
			return nil
		}
	})

	if err := eg.Wait(); err != nil {
		if errs.Is(err, ErrInterrupt) {
			slog.Error(
				fmt.Sprintf("the batch `%s` has been interrupted", batch.Name()),
				slog.Any("error", err),
			)
		} else {
			slog.Error(
				fmt.Sprintf("the batch `%s` failed", batch.Name()),
				slog.Any("error", err),
			)
		}
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("the batch `%s` finished successfully.", batch.Name()))
}

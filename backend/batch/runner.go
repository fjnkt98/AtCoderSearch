package batch

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/goark/errs"
	"golang.org/x/sync/errgroup"
)

var ErrInterrupt = errs.New("interrupted")

type Batch interface {
	Run(ctx context.Context) error
	Name() string
	Config() any
}

type Done struct{}

func RunBatch(batch Batch) {
	slog.Info("Start batch", slog.String("name", batch.Name()), slog.Any("config", batch.Config()))

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
				"the batch has been interrupted",
				slog.String("name", batch.Name()),
				slog.Any("config", batch.Config()),
				slog.Any("error", err),
			)
		} else {
			slog.Error(
				"the batch failed",
				slog.String("name", batch.Name()),
				slog.Any("config", batch.Config()),
				slog.Any("error", err),
			)
		}
		os.Exit(1)
	}
	slog.Info("the batch finished successfully.", slog.String("name", batch.Name()), slog.Any("config", batch.Config()))
}

package cmd

import (
	pb "fjnkt98/atcodersearch/grpc/atcodersearch/v1"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/searchers"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/meilisearch/meilisearch-go"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type Problem struct {
}

func NewServeCmd() *cli.Command {
	return &cli.Command{
		Name: "serve",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Value:   8000,
				EnvVars: []string{"BACKEND_PORT"},
			},
		},
		Action: func(c *cli.Context) error {
			port := c.Int("port")
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				return fmt.Errorf("listen tcp at port %d: %w", port, err)
			}

			pool, err := repository.NewPool(c.Context, c.String("database-url"))
			if err != nil {
				return fmt.Errorf("new pool: %w", err)
			}

			client, err := meilisearch.Connect(c.String("engine-url"), meilisearch.WithAPIKey(c.String("engine-master-key")))
			if err != nil {
				return fmt.Errorf("new connection: %w", err)
			}

			ctx, stop := signal.NotifyContext(c.Context, os.Interrupt)
			defer stop()

			s := grpc.NewServer(
				grpc.UnaryInterceptor(
					recovery.UnaryServerInterceptor(),
				),
			)

			pb.RegisterSearchServiceServer(s, searchers.NewSearcher(client, pool))

			go func() {
				slog.LogAttrs(ctx, slog.LevelInfo, "start grpc server", slog.Int("port", port))
				s.Serve(listener)
			}()

			<-ctx.Done()
			slog.LogAttrs(ctx, slog.LevelInfo, "stop grpc server")
			s.GracefulStop()

			return nil
		},
	}
}

package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	pb "fjnkt98/atcodersearch/grpc"
	"fjnkt98/atcodersearch/server"
)

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

			ctx, stop := signal.NotifyContext(c.Context, os.Interrupt)
			defer stop()

			s := grpc.NewServer()

			pb.RegisterProblemServiceServer(s, server.NewProblemSearcher())

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

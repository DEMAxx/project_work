package main

import (
	"context"
	"fmt"
	lrucache "github.com/DEMAxx/project_work/internal/lru_cache"
	internalhttp "github.com/DEMAxx/project_work/internal/server/http"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DEMAxx/project_work/pkg/config"
	"github.com/DEMAxx/project_work/pkg/logger"
)

func main() {
	const op = "cmd.previewer.main"

	cnf := config.MustLoad()

	logs := logger.MustSetupLogger(config.AppName, cnf.Env, cnf.Debug || cnf.Local, cnf.LogLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = logs.WithContext(ctx)

	cache := lrucache.NewCache(cnf.Capability)

	server := internalhttp.NewServer(
		&logs,
		net.JoinHostPort(cnf.Server.Host, cnf.Server.Port),
		cache,
		cnf,
	)

	ctx, cancel = signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if err := server.Start(ctx); err != nil {
		logs.Error().Msg(fmt.Sprintf("failed to start http server: %s", err.Error()))
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	logs.Info().Msg("calendar is running...")

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logs.Error().Msg(fmt.Sprintf("failed to stop http server: %s", err.Error()))
		}
	}()
}

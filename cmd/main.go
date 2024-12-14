package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGPIPE,
	)
	defer stop()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.Error("Err", slog.Any("err", err))
	}
}

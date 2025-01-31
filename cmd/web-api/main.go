package main

import (
	"context"
	"flag"
	"log/slog"

	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mohammadne/ice-global/internal/api/http"
	"github.com/mohammadne/ice-global/internal/config"
	"github.com/mohammadne/ice-global/pkg/mysql"
)

func main() {
	httpPort := flag.Int("http-port", 8088, "The server port which handles http requests (default: 8088)")
	flag.Parse() // Parse the command-line flags

	cfg, err := config.Load(true)
	if err != nil {
		panic(err)
	}

	mysql, err := mysql.Open(cfg.Mysql, "")
	if err != nil {
		slog.Error(`error connecting to mysql database`, `Err`, err)
		os.Exit(1)
	}
	_ = mysql

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go http.New().Serve(ctx, &wg, *httpPort)

	<-ctx.Done()
	wg.Wait()
	slog.Warn("interruption signal recieved, gracefully shutdown the server")
}

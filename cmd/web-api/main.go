package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"os/signal"
	"sync"
	"syscall"

	"github.com/mohammadne/ice-global/cmd"
	"github.com/mohammadne/ice-global/internal/api/http"
	"github.com/mohammadne/ice-global/internal/config"
	"github.com/mohammadne/ice-global/internal/services"
	"github.com/mohammadne/ice-global/internal/storage"
	"github.com/mohammadne/ice-global/pkg/mysql"
	"github.com/mohammadne/ice-global/pkg/redis"
)

func main() {
	httpPort := flag.Int("http-port", 8088, "The server port which handles http requests (default: 8088)")
	flag.Parse() // Parse the command-line flags

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo})))
	cmd.BuildInfo()

	cfg, err := config.LoadDefaults(true, "")
	if err != nil {
		panic(err)
	}

	mysql, err := mysql.Open(cfg.Mysql)
	if err != nil {
		slog.Error(`error connecting to mysql database`, `Err`, err)
		os.Exit(1)
	}

	redis, err := redis.Open(cfg.Redis)
	if err != nil {
		slog.Error(`error connecting to redis database`, `Err`, err)
		os.Exit(1)
	}

	// storages
	cartItemsStorage := storage.NewCartItems(mysql)
	cartsStorage := storage.NewCarts(mysql)
	itemsStorage := storage.NewItems(mysql, redis)

	// services
	cartsService := services.NewCarts(cartItemsStorage, cartsStorage, itemsStorage)
	itemsService := services.NewItems(itemsStorage)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go http.New(cartsService, itemsService).Serve(ctx, &wg, *httpPort)

	<-ctx.Done()
	wg.Wait()
	slog.Warn("interruption signal recieved, gracefully shutdown the server")
}

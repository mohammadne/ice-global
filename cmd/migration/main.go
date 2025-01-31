package main

import (
	"flag"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mohammadne/ice-global/internal/config"
	"github.com/mohammadne/ice-global/pkg/mysql"
)

func main() {
	migration := flag.String("migration", "hacks/schemas", "The migration directory (default: hacks/schemas)")
	direction := flag.String("direction", "", "Either 'UP' or 'DOWN'")
	flag.Parse() // Parse the command-line flags

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo})))
	slog.Info(`Go`, `Version`, runtime.Version(), `OS`, runtime.GOOS, `ARCH`, runtime.GOARCH, `now`, time.Now(), `Local`, time.Local)

	cfg, err := config.LoadDefaults(true, "")
	if err != nil {
		slog.Error(`error loading configs`, `Err`, err)
		os.Exit(1)
	}

	db, err := mysql.Open(cfg.Mysql, "file://"+*migration)
	if err != nil {
		slog.Error(`error connecting to mysql database`, `Err`, err)
		os.Exit(1)
	}

	migrateDirection := mysql.MigrateDirection(strings.ToUpper(*direction))
	err = db.Migrate(migrateDirection)
	if err != nil {
		slog.Error(`error migrating mysql database`, `Err`, err)
		os.Exit(1)
	}

	slog.Info(`database has been migrated`)
}

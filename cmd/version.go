package cmd

import (
	"fmt"
	"log/slog"
	"runtime"
)

// Default build-time variable.
// These values are overridden via ldflags
var (
	Version    = "unknown-version"
	GitCommit  = "unknown-commit"
	BuildTime  = "unknown-buildtime"
	APIVersion = "v0.1.0"
)

func BuildInfo() {
	slog.Info(`Build Information`,
		slog.String(`Version`, Version),
		slog.String(`API Version`, APIVersion),
		slog.String(`Go Version`, runtime.Version()),
		slog.String(`Git Commit`, GitCommit),
		slog.String(`Built At`, BuildTime),
		slog.String(`OS/ARCH`, fmt.Sprintf("OS/Arch:\t %s/%s\n", runtime.GOOS, runtime.GOARCH)),
	)
}

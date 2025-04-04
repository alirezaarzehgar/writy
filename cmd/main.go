package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/alirezaarzehgar/writy/internal/server"
	"github.com/alirezaarzehgar/writy/internal/writy"
)

func main() {
	dbPath := flag.String("db", writy.DefaultStoragePath, "database path for indexing and storage")
	runningAddr := flag.String("addr", ":8000", "running address e.g: localhost:8000, :3000. etc")
	reflecEnabled := flag.Bool("reflection", false, "enabled gRPC reflection for testing")
	logLevel := flag.String("leveler", "error", "log levels: error, warn, info, debug")
	flag.Parse()

	level := slog.LevelError
	switch *logLevel {
	case "warn", "4":
		level = slog.LevelWarn
	case "info", "0", "6":
		level = slog.LevelInfo
	case "debug", "-4", "7":
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	conf := server.ServerConfig{
		DbPath:            *dbPath,
		RunningAddr:       *runningAddr,
		ReflectionEnabled: *reflecEnabled,
	}

	slog.Debug("start gRPC server", "server config", conf, "leveler", *logLevel)
	err := server.Start(conf)
	if err != nil {
		slog.Error("failed to start server", "error", err)
	}

	fmt.Println(os.Args)
}

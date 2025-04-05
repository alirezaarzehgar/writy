package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/alirezaarzehgar/writy/internal/balancer"
	"github.com/alirezaarzehgar/writy/internal/server"
	"github.com/alirezaarzehgar/writy/internal/writy"
)

func main() {
	var slaves, masters balancer.StringArray
	flag.Var(&slaves, "slave", "list of slaves url for loadbalancing")
	flag.Var(&masters, "master", "list of slaves url for loadbalancing")
	runningAddr := flag.String("addr", ":8000", "running address e.g: localhost:8000, :3000. etc")
	reflecEnabled := flag.Bool("reflection", false, "enabled gRPC reflection for testing")
	logLevel := flag.String("leveler", "error", "log levels: error, warn, info, debug")
	isLoadbalancer := flag.Bool("balancer", false, "enable balancer to run a loadbalancer")
	dbPath := flag.String("db", writy.DefaultStoragePath, "database path for indexing and storage")
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

	if *isLoadbalancer {
		conf := balancer.ServerConfig{
			RunningAddr:       *runningAddr,
			ReflectionEnabled: *reflecEnabled,
			Slaves:            slaves,
			Masters:           masters,
		}

		slog.Debug("start gRPC server", "server config", conf, "leveler", *logLevel)
		err := balancer.Start(conf)
		if err != nil {
			slog.Error("failed to start loadbalancer", "error", err)
		}
	} else {
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
	}
}

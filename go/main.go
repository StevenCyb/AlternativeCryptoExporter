package main

import (
	slog "github.com/StevenCyb/SimpleLogging"
)

func main() {
	slog.Initialize(
		slog.LogLevel(logLevel),
		slog.FORMAT_JSON,
		slog.LogFiles{Path: ""},
	)

	slog.Info(slog.EntryEvent("Start watchers..."))
	err := setupWatcher()
	if err != nil {
		slog.Fatal(slog.Entry{"event": "Startup failed", "error": err.Error()})
	}

	slog.Info(slog.EntryEvent("Run metrics server..."))
	startMetricsServer()
}

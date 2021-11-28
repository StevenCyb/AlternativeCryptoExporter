package main

import (
	"AlternativeCryptoExporter/config"
	"AlternativeCryptoExporter/crypto"

	slog "github.com/StevenCyb/SimpleLogging"
)

func main() {
	slog.Initialize(
		slog.LogLevel(config.LogLevel),
		slog.FORMAT_JSON,
		slog.LogFiles{Path: ""},
	)

	slog.Info(slog.Entry{"event": "Starting..."})
	cryptoAPI := crypto.API{}

	slog.Info(slog.Entry{"event": "Retrieve list of supported currency..."})
	err := cryptoAPI.UpdateSupportedCryptoList()
	if err != nil {
		slog.Fatal(slog.Entry{"event": "Retrieve list of supported currency failed", "error": err.Error()})
	}

	slog.Info(slog.Entry{"event": "Starting watcher..."})
	err = cryptoAPI.StartWatcher()
	if err != nil {
		slog.Fatal(slog.Entry{"event": "Failed to start watcher", "error": err.Error()})
	}
}

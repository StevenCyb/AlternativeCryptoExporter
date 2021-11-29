package main

import (
	"AlternativeCryptoExporter/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	slog "github.com/StevenCyb/SimpleLogging"
)

func setupWatcher() error {
	slog.Debug(slog.EntryEvent("Requesting list of supported crypto curracy"))
	resp, err := http.Get("https://api.alternative.me/v2/listings/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	slog.Debug(slog.EntryEvent("Read body"))
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	currencyListings := model.CurrencyListings{}
	slog.Debug(slog.EntryEvent("Unmarshal body"))
	err = json.Unmarshal(bodyBytes, &currencyListings)

	slog.Debug(slog.EntryEvent("Create watcher"))
	duplicateList := []string{}
	for _, watchCurrency := range currencyWatchlist {
		found := false

		for _, supportedCurrency := range currencyListings.Data {
			if strings.EqualFold(watchCurrency, supportedCurrency.Name) || strings.EqualFold(watchCurrency, supportedCurrency.Symbol) {
				found = true
				duplicate := false

				for _, id := range duplicateList {
					if id == supportedCurrency.ID {
						duplicate = true
						break
					}
				}
				if duplicate {
					slog.Warning(slog.EntryEvent(fmt.Sprintf("Skip unsupported currency '%s'", supportedCurrency.Name)))
					break
				}

				duplicateList = append(duplicateList, supportedCurrency.ID)

				runTicker(supportedCurrency)
				break
			}
		}

		if !found {
			return fmt.Errorf("unsupported currency '%s'", watchCurrency)
		}
	}

	return err
}

func runTicker(currency model.CurrencyListingData) {
	// Interval is statically set to 5 min since the API
	// refreshed every 5 min
	duration := time.Duration(5 * time.Second)
	go func() {
		for range time.Tick(duration) {
			slog.Info(slog.Entry{
				"even":            "Update currency value",
				"currency_name":   currency.Name,
				"currency_symbol": currency.Symbol,
				"currency_id":     currency.ID,
			})
		}
	}()
}

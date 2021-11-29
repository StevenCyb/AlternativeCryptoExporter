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
	quotes = strings.ToUpper(quotes)

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

	cryptoCurrencyListings := model.CryptoCurrencyListings{}
	slog.Debug(slog.EntryEvent("Unmarshal body"))
	err = json.Unmarshal(bodyBytes, &cryptoCurrencyListings)

	slog.Debug(slog.EntryEvent("Create watcher"))
	duplicateList := []string{}
	for _, watchCurrency := range currencyWatchlist {
		found := false

		for _, supportedCurrency := range cryptoCurrencyListings.Data {
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

func getLogEntryWithCryptoCurrency(cryptoCurrency model.CryptoCurrencyListingData, event string, err error) slog.Entry {
	entry := slog.Entry{
		"even":                   event,
		"crypto_currency_name":   cryptoCurrency.Name,
		"crypto_currency_symbol": cryptoCurrency.Symbol,
		"crypto_currency_id":     cryptoCurrency.ID,
	}
	if err != nil {
		entry["error"] = err.Error()
	}

	return entry
}

func runTicker(cryptoCurrency model.CryptoCurrencyListingData) {
	// Interval is statically set to 5 min since the API
	// refreshed every 5 min
	duration := time.Duration(5 * time.Minute)

	UpdateMetrics(cryptoCurrency)
	go func() {
		for range time.Tick(duration) {
			UpdateMetrics(cryptoCurrency)
		}
	}()
}

func UpdateMetrics(cryptoCurrency model.CryptoCurrencyListingData) {
	slog.Debug(
		getLogEntryWithCryptoCurrency(cryptoCurrency,
			"Requesting crypto currency data", nil))
	resp, err := http.Get(fmt.Sprintf(
		"https://api.alternative.me/v2/ticker/%s/?structure=array&convert=%s",
		cryptoCurrency.ID, quotes))
	if err != nil {
		slog.Error(
			getLogEntryWithCryptoCurrency(cryptoCurrency,
				"Requesting crypto currency data failed", err))
		return
	}

	slog.Debug(
		getLogEntryWithCryptoCurrency(cryptoCurrency,
			"Read body", nil))
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(
			getLogEntryWithCryptoCurrency(cryptoCurrency,
				"Reading response data failed", err))
		return
	}
	resp.Body.Close()

	cryptoCurrencyData := model.CryptoCurrency{}
	slog.Debug(
		getLogEntryWithCryptoCurrency(cryptoCurrency,
			"Unmarshal body", nil))
	err = json.Unmarshal(bodyBytes, &cryptoCurrencyData)
	if err != nil {
		slog.Error(
			getLogEntryWithCryptoCurrency(cryptoCurrency,
				"Unmarshal response data failed", err))
		return
	}

	if len(cryptoCurrencyData.Data) != 1 {
		slog.Warning(
			getLogEntryWithCryptoCurrency(cryptoCurrency,
				"Response has more or les then one data entry", nil))
		return
	}

	CryptoCurrencyCirculatingSupplyGauge.
		WithLabelValues(cryptoCurrency.ID, cryptoCurrency.Name, cryptoCurrency.Symbol, quotes).
		Set(cryptoCurrencyData.Data[0].CirculatingSupply)

	CryptoCurrencyTotalSupplyGauge.
		WithLabelValues(cryptoCurrency.ID, cryptoCurrency.Name, cryptoCurrency.Symbol, quotes).
		Set(cryptoCurrencyData.Data[0].TotalSupply)

	CryptoCurrencyMaxSupplyGauge.
		WithLabelValues(cryptoCurrency.ID, cryptoCurrency.Name, cryptoCurrency.Symbol, quotes).
		Set(cryptoCurrencyData.Data[0].MaxSupply)

	if entry, ok := cryptoCurrencyData.Data[0].Quotes[quotes]; ok {
		CryptoCurrencyPriceGauge.
			WithLabelValues(cryptoCurrency.ID, cryptoCurrency.Name, cryptoCurrency.Symbol, quotes).
			Set(entry.Price)

		CryptoCurrencyVolume24hGauge.
			WithLabelValues(cryptoCurrency.ID, cryptoCurrency.Name, cryptoCurrency.Symbol, quotes).
			Set(entry.Volume24H)

		CryptoCurrencyMarketCapGauge.
			WithLabelValues(cryptoCurrency.ID, cryptoCurrency.Name, cryptoCurrency.Symbol, quotes).
			Set(entry.MarketCap)

		CryptoCurrencyPriceChangePercentageGauge.
			WithLabelValues(cryptoCurrency.ID, cryptoCurrency.Name, cryptoCurrency.Symbol, quotes, "1h").
			Set(entry.PercentChange1H)
		CryptoCurrencyPriceChangePercentageGauge.
			WithLabelValues(cryptoCurrency.ID, cryptoCurrency.Name, cryptoCurrency.Symbol, quotes, "24h").
			Set(entry.PercentChange24H)
		CryptoCurrencyPriceChangePercentageGauge.
			WithLabelValues(cryptoCurrency.ID, cryptoCurrency.Name, cryptoCurrency.Symbol, quotes, "7d").
			Set(entry.PercentChange7D)
	} else {
		slog.Warning(
			getLogEntryWithCryptoCurrency(cryptoCurrency,
				fmt.Sprintf("Response not contains quotes '%s'", quotes), nil))
		return
	}

	slog.Info(
		getLogEntryWithCryptoCurrency(cryptoCurrency,
			"Crypto currency updated", nil))
}

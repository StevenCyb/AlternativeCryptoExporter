package crypto

import (
	"AlternativeCryptoExporter/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	slog "github.com/StevenCyb/SimpleLogging"
)

type API struct {
	supportedCurrency SupportedCurrency
}

func (api *API) UpdateSupportedCryptoList() error {
	slog.Debug(slog.Entry{"event": "Requesting list"})
	resp, err := http.Get("https://api.alternative.me/v2/listings/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	slog.Debug(slog.Entry{"event": "Read body"})
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	slog.Debug(slog.Entry{"event": "Unmarshal body"})
	err = json.Unmarshal(bodyBytes, &api.supportedCurrency)

	return err
}

func (api *API) StartWatcher() error {
	slog.Debug(slog.Entry{"event": "Check and run watchlist"})
	duplicateList := []string{}

	for _, watchCurrency := range config.CurrencyWatchlist {
		found := false

		for _, supportedCurrency := range api.supportedCurrency.Data {
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
					break
				}

				duplicateList = append(duplicateList, supportedCurrency.ID)

				fmt.Printf("TODO: Watch %v\n", supportedCurrency)
				break
			}
		}

		if !found {
			return fmt.Errorf("unsupported currency '%s'", watchCurrency)
		}
	}

	return nil
}

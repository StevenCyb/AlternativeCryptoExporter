package main

import (
	"AlternativeCryptoExporter/model"
	"flag"
)

var (
	logLevel          int
	listen            string
	quotes            string
	currencyWatchlist model.StringArrayFlags
)

func init() {
	flag.IntVar(&logLevel, "log-level", 4, "Log level {'0' NONE, '1' FATAL, '2' ERROR, '3' WARNING, '4' INFO, '5' DEBUG}")
	flag.StringVar(&listen, "listen", ":8080", "Set listen endpoint e.g ':8080'")

	flag.StringVar(&quotes, "quotes", "USD", "Currency to display crypto quotes e.g. {USD', 'EUR', 'GBP', 'RUB', 'JPY', 'CAD', 'KRW', 'PLN', 'BTC', 'ETH', 'XRP', 'LTC}")
	flag.Var(&currencyWatchlist, "watch", "Define currency to watch e.g. {'BTC', 'LTC'}")

	flag.Parse()
}

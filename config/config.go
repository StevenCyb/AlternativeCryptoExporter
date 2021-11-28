package config

import (
	"flag"
)

type StringArrayFlags []string

func (saf *StringArrayFlags) String() string {
	out := "["
	for _, s := range *saf {
		out = out + s + ","
	}
	return out[:len(out)-1] + "]"
}

func (saf *StringArrayFlags) Set(value string) error {
	*saf = append(*saf, value)
	return nil
}

var (
	LogLevel          int
	ConvertCurrency   string
	CurrencyWatchlist StringArrayFlags
)

func init() {
	flag.IntVar(&LogLevel, "log-level", 4, "Log level {'0' NONE, '1' FATAL, '2' ERROR, '3' WARNING, '4' INFO, '5' DEBUG}")
	flag.StringVar(&ConvertCurrency, "convert-currency", "", "Currency to display crypto value e.g. {USD', 'EUR', 'GBP', 'RUB', 'JPY', 'CAD', 'KRW', 'PLN', 'BTC', 'ETH', 'XRP', 'LTC}")
	flag.Var(&CurrencyWatchlist, "watch", "Define currency to watch e.g. {'BTC', 'LTC'}")

	flag.Parse()
}

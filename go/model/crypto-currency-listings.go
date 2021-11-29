package model

type CryptoCurrencyListings struct {
	Data []CryptoCurrencyListingData `json:"data"`
}

type CryptoCurrencyListingData struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

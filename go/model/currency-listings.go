package model

type CurrencyListings struct {
	Data []CurrencyListingData `json:"data"`
}

type CurrencyListingData struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

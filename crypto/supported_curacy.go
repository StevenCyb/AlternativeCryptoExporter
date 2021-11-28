package crypto

type SupportedCurrency struct {
	Data []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Watch  bool
	} `json:"data"`
}

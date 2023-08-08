package requests

type ExchangeRate struct {
	CurrencyTo          string `json:"CurrencyTo"`
	CurrencyFrom        string `json:"CurrencyFrom"`
	ValidityStartDate   string `json:"ValidityStartDate"`
	ValidityEndDate     string `json:"ValidityEndDate"`
	IsMarkedForDeletion *bool  `json:"IsMarkedForDeletion"`
}

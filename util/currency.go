package util

// Constants for all supported currency.
const (
	USD = "USD"
	CAD = "CAD"
	EUR = "EUR"
)

// IsSupportedCurrency returns true if the currency is supported, false otherwise.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, CAD, EUR:
		return true
	}
	return false
}

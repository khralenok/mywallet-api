package utilities

func ConvertToUSD(cents int) float64 {
	dollars := float64(cents) / 100

	return dollars
}

func ConvertToCents(dollars float64) int {
	cents := int(dollars * 100)

	return cents
}

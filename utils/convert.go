package utils

import (
	"fmt"
	"math"
)

func UnitUinUt64(n uint64) string {
	return UnitFloat(float64(n))
}

func UnitFloat(number float64) string {
	units := []string{"", "K", "M", "G"}
	base := 1024.0

	if number < base {
		return fmt.Sprintf("%.0f", number)
	}

	exp := int(math.Log(number) / math.Log(base))
	scaledNumber := number / math.Pow(base, float64(exp))
	unit := units[exp]

	return fmt.Sprintf("%.1f%s", scaledNumber, unit)
}

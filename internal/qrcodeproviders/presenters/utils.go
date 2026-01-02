package presenters

import (
	"fmt"
	"strconv"
)

func FormatDecimal(value float64) float64 {
	formattedValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return formattedValue
}

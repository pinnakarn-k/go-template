package formatter

import (
	"strings"

	"github.com/shopspring/decimal"
)

func FormatDecimal(value *decimal.Decimal) *string {
	if value == nil {
		return nil
	}

	formatted := addThousandsSeparator(value.StringFixed(2))
	return &formatted
}

func addThousandsSeparator(value string) string {
	parts := strings.SplitN(value, ".", 2)

	integerPart := parts[0]
	decimalPart := ""

	if len(parts) == 2 {
		decimalPart = "." + parts[1]
	}

	sign := ""
	if strings.HasPrefix(integerPart, "-") {
		sign = "-"
		integerPart = strings.TrimPrefix(integerPart, "-")
	}

	var builder strings.Builder

	firstGroupLength := len(integerPart) % 3
	if firstGroupLength == 0 {
		firstGroupLength = 3
	}

	builder.WriteString(integerPart[:firstGroupLength])

	for i := firstGroupLength; i < len(integerPart); i += 3 {
		builder.WriteByte(',')
		builder.WriteString(integerPart[i : i+3])
	}

	return sign + builder.String() + decimalPart
}

package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func ToString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		if v == "" || strings.EqualFold(v, "null") {
			return ""
		}
		return strings.TrimSpace(v)
	case float64:
		return strings.TrimSpace(strconv.FormatFloat(v, 'f', -1, 64))
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func ToPtrString(value any) *string {
	switch v := value.(type) {
	case nil:
		return nil
	default:
		s := ToString(v)
		if s == "" {
			return nil
		}
		return &s
	}
}

func ToDecimal(value any) (decimal.Decimal, error) {
	switch v := value.(type) {
	case nil:
		return decimal.Zero, nil
	case string:
		if v == "" || strings.EqualFold(v, "null") {
			return decimal.Zero, nil
		}
		return decimal.NewFromString(v)
	case float64:
		s := strings.TrimSpace(strconv.FormatFloat(v, 'f', -1, 64))
		if s == "" {
			return decimal.Zero, nil
		}
		return decimal.NewFromString(s)
	default:
		s := strings.TrimSpace(fmt.Sprint(v))
		if s == "" {
			return decimal.Zero, nil
		}
		return decimal.NewFromString(s)
	}
}

func ToPtrDecimal(value any) (*decimal.Decimal, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	default:
		d, err := ToDecimal(v)
		if err != nil {
			return nil, err
		}
		if d.Equal(decimal.Zero) {
			return nil, nil
		}
		return &d, nil

	}
}

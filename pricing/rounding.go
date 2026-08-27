package pricing

import (
	"errors"
	"math"
)

type RoundingMode string

const (
	RoundNearest RoundingMode = "nearest"
	RoundUp      RoundingMode = "up"
	RoundDown    RoundingMode = "down"
)

func ApplyPercent(cents int64, percent int64, mode RoundingMode) (int64, error) {
	if cents <= 0 || percent < 0 {
		return 0, errors.New("cents and percent must be valid")
	}
	if mode != RoundNearest && mode != RoundUp && mode != RoundDown {
		return 0, errors.New("unknown rounding mode")
	}
	raw := float64(cents) * (100 + float64(percent)) / 100
	switch mode {
	case RoundUp:
		return int64(math.Ceil(raw)), nil
	case RoundDown:
		return int64(math.Floor(raw)), nil
	default:
		return int64(math.Round(raw)), nil
	}
}

func Discount(cents int64, percent int64) (int64, error) {
	if percent < 0 || percent > 100 {
		return 0, errors.New("discount must be between zero and one hundred")
	}
	return cents - cents*percent/100, nil
}

func IsAffordable(cents, budget int64) bool {
	return cents > 0 && budget >= cents
}

func TaxInclusive(cents, basisPoints int64) (int64, error) {
	if cents <= 0 || basisPoints < 0 {
		return 0, errors.New("invalid tax input")
	}
	return cents + cents*basisPoints/10000, nil
}

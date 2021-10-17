package numb

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func ToDecimal(hex string, exp int32) decimal.Decimal {
	bi, ok := new(big.Int).SetString(hex, 0)
	if !ok {
		return decimal.NewFromInt(0)
	}
	return decimal.NewFromBigInt(bi, exp)
}

func ToGwei(hex string) (decimal.Decimal, error) {
	bi, ok := new(big.Int).SetString(hex, 0)
	if !ok {
		return decimal.Decimal{}, fmt.Errorf("converting to gwei")
	}
	return decimal.NewFromBigInt(bi, -9).Floor(), nil
}

func ToInt64(hex string) (int64, error) {
	value, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parsing timestamp")
	}
	return value, nil
}

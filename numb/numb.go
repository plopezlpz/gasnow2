package numb

import (
	"math/big"
	"strconv"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func ToDecimal(hex string, exp int32) decimal.Decimal {
	bi, ok := new(big.Int).SetString(hex, 0)
	if !ok {
		decimal.NewFromInt(0)
	}
	return decimal.NewFromBigInt(bi, exp)
}

func ToTimestamp(hex string) (int64, error) {
	value, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parsing timestamp")
	}
	return value, nil
}

package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type RawTransaction struct {
	SK                   string
	From                 *common.Address
	To                   *common.Address
	MethodId             [4]byte
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
	GasLimit             uint64
	Value                *big.Int
	Data                 []byte
}

package controller

import (
	"fmt"
	"math/big"
)

type RequestFaucet struct {
	Address     string `json:"address" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
	TokenSymbol string `json:"token_symbol" binding:"required"`
	ChainID     string `json:"chain_id" binding:"required"`
}

func ethToWei(ethAmount string) (*big.Int, error) {
	fmt.Println(ethAmount)

	// 将 ETH 转换为 Wei
	amountFloat, ok := new(big.Float).SetString(ethAmount)
	if !ok {
		return nil, fmt.Errorf("invalid amount format: %s", ethAmount)
	}

	// 1 ETH = 10^18 wei
	// 使用 big.Float 来处理乘法操作，保持精度
	weiFloat := new(big.Float).Mul(amountFloat, big.NewFloat(1e18)) // ETH -> wei

	// 转换为 big.Int，注意这里是通过 Int() 方法将 big.Float 转换为 big.Int
	weiInt, _ := weiFloat.Int(nil)

	// 这样 weiInt 会是一个精确的整数表示
	return weiInt, nil
}

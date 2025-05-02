package controller

import (
	"fmt"
	"math/big"
)

type FaucetRequest struct {
	Address string `json:"address" binding:"required"`
	Token   string `json:"token" binding:"required"`
}

type SignRequest struct {
	Code string `json:"code" binding:"required"`
}

type SignResponse struct {
	Token string `json:"token"`
}

type FaucetResponse struct {
	Address string `json:"address"`
	Tx      string `json:"tx"`
}

type AccessTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

// 定义响应的结构体
type AccessTokenResponse struct {
	Status int `json:"status"`
	Code   int `json:"code"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

// 定义数据部分的结构体
type UserData struct {
	Uid      uint   `json:"uid"`
	Avatar   string `json:"avatar"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Github   string `json:"github"`
}

// 定义顶层响应的结构体
type GetUserResponse struct {
	ID      string   `json:"id"`
	Status  int      `json:"status"`
	Code    int      `json:"code"`
	Data    UserData `json:"data"`
	Time    int64    `json:"time"`
	Message string   `json:"message"`
}

func ethToWei(ethAmount string) (*big.Int, error) {
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

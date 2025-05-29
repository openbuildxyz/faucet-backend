package model

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Address      string `json:"address" gorm:"not null"`      // 钱包地址
	Amount       string `json:"amount" gorm:"not null"`       // 领水的金额（单位：wei）
	TxHash       string `json:"tx_hash" gorm:"not null"`      // 交易哈希
	Status       string `json:"status" gorm:"not null"`       // 状态，可能是 "pending", "completed", "failed"
	TokenSymbol  string `json:"token_symbol" gorm:"not null"` // 代币符号，例如 ETH, USDT 等
	ChainType    string `json:"chain_type" gorm:"not null"`   // 链类型，例如 Ethereum, Binance Smart Chain 等
	ChainID      string `json:"chain_id"`                     // 链 ID，例如 Ethereum 主网为 1，BSC 为 56
	RpcURL       string `json:"rpc_url"`                      // 对应链的 RPC URL，例如 Infura、Alchemy 等 URL
	ErrorMessage string `json:"error_message"`                // 错误信息
	Uid          uint   `json:"uid"`
	Github       string `json:"github"`
}

func CreateTransaction(t *Transaction) error {
	if err := db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func GetTransactionByAddress(address, token string) (*Transaction, error) {
	var t Transaction
	if err := db.Where("address = ?", address).Where("token_symbol = ?", token).Last(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func GetTransactionByUid(uid uint, token string) (*Transaction, error) {
	var t Transaction
	if err := db.Where("uid = ?", uid).Where("token_symbol = ?", token).Last(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func GetTransactionByGithub(github, token string) (*Transaction, error) {
	var t Transaction
	if err := db.Where("github = ?", github).Where("token_symbol = ?", token).Last(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateWalletStatus 更新领水记录的状态
func UpdateWalletStatus(txhash, status string) error {
	var t Transaction
	if err := db.Where("tx_hash = ?", txhash).First(&t).Error; err != nil {
		return err
	}

	t.Status = status
	if err := db.Save(&t).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTransacton(t *Transaction) error {
	if err := db.Save(t).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTransaction(txhash string) error {
	var t Transaction
	if err := db.Where("tx_hash = ?", t).First(&t).Error; err != nil {
		return err
	}
	if err := db.Delete(&t).Error; err != nil {
		return err
	}
	return nil
}

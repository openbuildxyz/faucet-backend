package model

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Address      string    `json:"address" gorm:"not null"`      // 钱包地址
	RequestedAt  time.Time `json:"requested_at" gorm:"not null"` // 请求领水的时间
	Amount       string    `json:"amount" gorm:"not null"`       // 领水的金额（单位：wei）
	TxHash       string    `json:"tx_hash" gorm:"not null"`      // 交易哈希
	Status       string    `json:"status" gorm:"not null"`       // 状态，可能是 "pending", "completed", "failed"
	TokenSymbol  string    `json:"token_symbol" gorm:"not null"` // 代币符号，例如 ETH, USDT 等
	ChainType    string    `json:"chain_type" gorm:"not null"`   // 链类型，例如 Ethereum, Binance Smart Chain 等
	ChainID      string    `json:"chain_id" gorm:"not null"`     // 链 ID，例如 Ethereum 主网为 1，BSC 为 56
	RpcURL       string    `json:"rpc_url" gorm:"not null"`      // 对应链的 RPC URL，例如 Infura、Alchemy 等 URL
	ErrorMessage string    `json:"error_message"`                // 错误信息
}

func CreateWallet(wallet *Wallet) error {
	if err := db.Create(wallet).Error; err != nil {
		return err
	}
	return nil
}

// GetWalletByAddress 根据地址查找领水记录
func GetWalletByAddress(address string) (*Wallet, error) {
	var wallet Wallet
	if err := db.Where("address = ?", address).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetWalletsByStatus 根据状态查找领水记录
func GetWalletsByStatus(status string) ([]Wallet, error) {
	var wallets []Wallet
	if err := db.Where("status = ?", status).Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

// UpdateWalletStatus 更新领水记录的状态
func UpdateWalletStatus(address, status string) error {
	var wallet Wallet
	if err := db.Where("address = ?", address).First(&wallet).Error; err != nil {
		return err
	}

	wallet.Status = status
	if err := db.Save(&wallet).Error; err != nil {
		return err
	}
	return nil
}

// UpdateWallet 更新领水记录
func UpdateWallet(wallet *Wallet) error {
	if err := db.Save(wallet).Error; err != nil {
		return err
	}
	return nil
}

// DeleteWallet 根据地址删除领水记录
func DeleteWallet(address string) error {
	var wallet Wallet
	if err := db.Where("address = ?", address).First(&wallet).Error; err != nil {
		return err
	}
	if err := db.Delete(&wallet).Error; err != nil {
		return err
	}
	return nil
}

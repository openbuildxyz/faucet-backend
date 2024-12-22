package controller

import (
	"faucet/chain"
	"faucet/logger"
	"faucet/model"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

func HandleFaucet(c *gin.Context) {
	// 请求数据结构
	var req RequestFaucet
	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Errorf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 验证地址是否有效
	if !common.IsHexAddress(req.Address) {
		logger.Log.Errorf("Invalid Ethereum address: %s", req.Address)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Ethereum address"})
		return
	}

	// 生成新的钱包记录
	newWallet := &model.Wallet{
		Address:     req.Address,
		RequestedAt: time.Now(),
		Amount:      req.Amount,
		TxHash:      "",
		Status:      "pending", // 设置初始状态为 "pending"
		TokenSymbol: req.TokenSymbol,
		// ChainType:   req.ChainType,
		// ChainID:     req.ChainID,
		// RpcURL:      req.RpcURL,
	}

	// 保存记录到数据库
	if err := model.CreateWallet(newWallet); err != nil {
		logger.Log.Errorf("Failed to create wallet record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save wallet record"})
		return
	}

	// 假设你有一个 sendTransaction 方法，调用它来发送交易
	err := sendTransaction(req.Address, req.Amount)
	if err != nil {
		// 更新状态为 failed，并记录错误信息
		newWallet.Status = "failed"
		newWallet.ErrorMessage = err.Error()
		model.UpdateWallet(newWallet)

		logger.Log.Errorf("Failed to send transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send transaction"})
		return
	}

	// 更新状态为 success
	newWallet.Status = "success"
	if err := model.UpdateWallet(newWallet); err != nil {
		logger.Log.Errorf("Failed to update wallet status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet status"})
		return
	}

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "Faucet transaction sent successfully"})
}

// sendTransaction 发送代币或 ETH 的交易
func sendTransaction(address string, amount string) error {
	// TODO: 转换失去精度
	weiAmount, err := ethToWei(amount)
	if err != nil {
		return err
	}

	tx, err := chain.Transfer(address, weiAmount)
	if err != nil {
		return err
	}
	fmt.Println(tx)

	return nil
}

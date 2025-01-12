package controller

import (
	"faucet/chain"
	"faucet/logger"
	"faucet/model"
	"faucet/utils"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

func HandleFaucet(c *gin.Context) {
	var req FaucetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Errorf("Invalid request: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !common.IsHexAddress(req.Address) {
		logger.Log.Errorf("Invalid address: %s", req.Address)
		utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid address: %s", req.Address), nil)
		return
	}

	tx := &model.Transaction{
		Address:     req.Address,
		Amount:      req.Amount,
		TxHash:      "",
		Status:      "pending", // 设置初始状态为 "pending"
		TokenSymbol: req.TokenSymbol,
		// ChainType:   req.ChainType,
		// ChainID:     req.ChainID,
		// RpcURL:      req.RpcURL,
	}

	if err := model.CreateTransaction(tx); err != nil {
		logger.Log.Errorf("Failed to create transaction record: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	txhash, err := sendTransaction(req.Address, req.Amount)
	if err != nil {
		tx.Status = "failed"
		tx.ErrorMessage = err.Error()
		if err := model.UpdateTransacton(tx); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		logger.Log.Errorf("Failed to send transaction: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// TODO: check status(goroutines)

	var faucetResp FaucetResponse
	faucetResp.Address = req.Address
	faucetResp.Tx = txhash

	utils.SuccessResponse(c, http.StatusOK, "success", faucetResp)
}

func sendTransaction(address string, amount string) (string, error) {
	// TODO:
	weiAmount, err := ethToWei(amount)
	if err != nil {
		return "", err
	}

	tx, err := chain.Transfer(address, weiAmount)
	return tx, err
}

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
	"github.com/spf13/viper"
)

func HandleFaucet(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		logger.Log.Errorf("Invalid request: %v", "no token")
		utils.ErrorResponse(c, http.StatusUnauthorized, "", nil)
		return
	}

	user, err := model.GetUserByToken(authHeader)
	if err != nil {
		logger.Log.Errorf("Invalid request: %v", "no token")
		utils.ErrorResponse(c, http.StatusUnauthorized, "", nil)
		return
	}

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

	if req.ChainID != "20143" || req.TokenSymbol != "DMON" {
		logger.Log.Errorf("Invalid token info: chainid: %s, token: %s", req.ChainID, req.TokenSymbol)
		utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid token info %s %s", req.ChainID, req.TokenSymbol), nil)
		return
	}

	amountLimit := viper.GetString("monad.amount")
	if req.Amount != amountLimit {
		logger.Log.Errorf("only claim 1 DMON at a time, %s", req.Amount)
		utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("only claim %s DMON at a time", amountLimit), nil)
		return
	}

	wallet, err := model.GetTransactionByAddress(req.Address)
	if err == nil {
		if utils.IsWithinLast24Hours(wallet.CreatedAt) {
			logger.Log.Errorf("This wallet %s has already made a request. Please try again later.", req.Address)
			utils.ErrorResponse(c, http.StatusBadRequest, "This wallet has already made a request. Please try again later.", nil)
			return
		}
	}

	u, err := model.GetTransactionByUid(user.Uid)
	if err == nil {
		if utils.IsWithinLast24Hours(u.CreatedAt) {
			logger.Log.Errorf("This user %d has already made a request. Please try again later.", u.Uid)
			utils.ErrorResponse(c, http.StatusBadRequest, "You has already made a request. Please try again later.", nil)
			return
		}
	}

	tx := &model.Transaction{
		Address:     req.Address,
		Amount:      req.Amount,
		TxHash:      "",
		Status:      "pending", // 设置初始状态为 "pending"
		TokenSymbol: req.TokenSymbol,
		// ChainType:   req.ChainType,
		ChainID: req.ChainID,
		Uid:     user.Uid,
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

package controller

import (
	"errors"
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
	oauthToken, exists := c.Get("oauth_token")
	if !exists {
		logger.Log.Errorf("Invalid request: %v", "no token")
		utils.ErrorResponse(c, http.StatusUnauthorized, "Please log in to continue!", nil)
		return
	}

	oToken, ok := oauthToken.(string)
	if !ok {
		logger.Log.Errorf("Invalid request: %v", "no token")
		utils.ErrorResponse(c, http.StatusUnauthorized, "Please log in to continue!", nil)
	}

	user, err := model.GetUserByToken(oToken)
	if err != nil {
		logger.Log.Errorf("Invalid request: %v", "no token")
		utils.ErrorResponse(c, http.StatusUnauthorized, "Please log in to continue!", nil)
		return
	}

	var req FaucetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Errorf("Invalid request: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request. Please try again later.", nil)
		return
	}

	if !common.IsHexAddress(req.Address) {
		logger.Log.Errorf("Invalid address: %s", req.Address)
		utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid address: %s", req.Address), nil)
		return
	}

	if req.Token != "MON" && req.Token != "0G" && req.Token != "CAMP" {
		logger.Log.Errorf("Invalid token: %s", req.Token)
		utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid token: %s", req.Token), nil)
		return
	}

	wallet, err := model.GetTransactionByAddress(req.Address, req.Token)
	if err == nil {
		if utils.IsWithinLast24Hours(wallet.CreatedAt) {
			logger.Log.Errorf("This wallet %s has already made a request. Please try again later.", req.Address)
			utils.ErrorResponse(c, http.StatusBadRequest, "You has already made a request in 24 hours. Please try again later.", nil)
			return
		}
	}

	u, err := model.GetTransactionByUid(user.Uid, req.Token)
	if err == nil {
		if utils.IsWithinLast24Hours(u.CreatedAt) {
			logger.Log.Errorf("This user %d has already made a request. Please try again later.", u.Uid)
			utils.ErrorResponse(c, http.StatusBadRequest, "You has already made a request in 24 hours. Please try again later.", nil)
			return
		}
	}

	if user.Github == "" {
		logger.Log.Errorf("Please bind your GitHub in OpenBuiild first, %d, %v", u.Uid, *user)
		utils.ErrorResponse(c, http.StatusBadRequest, "Please bind your GitHub in OpenBuiild first", nil)
		return
	}

	g, err := model.GetTransactionByGithub(user.Github, req.Token)
	if err == nil {
		if utils.IsWithinLast24Hours(g.CreatedAt) {
			logger.Log.Errorf("This user %d, %s has already made a request. Please try again later.", g.Uid, g.Github)
			utils.ErrorResponse(c, http.StatusBadRequest, "You has already made a request in 24 hours. Please try again later.", nil)
			return
		}
	}

	var amount string
	amount, err = RequestGitRank(user.Github)
	if err != nil {
		logger.Log.Errorf("RequestGitRank error, %s", err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	// if req.Token == "MON" {
	// 	amount, err = RequestGitRankWithToken(user.Github, req.Token)
	// 	if err != nil {
	// 		logger.Log.Errorf("RequestGitRank error, %s", err.Error())
	// 		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
	// 		return
	// 	}
	// }

	var chain string
	if req.Token == "MON" {
		chain = "10143"
	}
	if req.Token == "0G" {
		chain = "16601"
	}
	if req.Token == "CAMP" {
		chain = "123420001114"
	}

	tx := &model.Transaction{
		Address:     req.Address,
		Amount:      amount,
		TxHash:      "",
		Status:      "pending", // 设置初始状态为 "pending"
		TokenSymbol: req.Token,
		ChainType:   "evm",
		ChainID:     chain,
		Uid:         user.Uid,
		Github:      user.Github,
		// RpcURL:      req.RpcURL,
	}

	if err := model.CreateTransaction(tx); err != nil {
		logger.Log.Errorf("Failed to create transaction record: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	txhash, err := sendTransaction(req.Address, req.Token, amount)
	if err != nil {
		tx.Status = "failed"
		tx.ErrorMessage = err.Error()
		if err := model.UpdateTransacton(tx); err != nil {
			logger.Log.Errorf("Failed to update transaction record: %s, %s", req.Address, err.Error())
			utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
			return
		}

		logger.Log.Errorf("Failed to send transaction: %s, %s", req.Address, err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	// TODO: check status(goroutines)

	var faucetResp FaucetResponse
	faucetResp.Address = req.Address
	faucetResp.Tx = txhash

	utils.SuccessResponse(c, http.StatusOK, "success", faucetResp)
}

func sendTransaction(address, token, amount string) (string, error) {
	// TODO:
	weiAmount, err := ethToWei(amount)
	if err != nil {
		return "", err
	}

	tx, err := chain.Transfer(address, token, weiAmount)
	return tx, err
}

func RequestGitRank(github string) (string, error) {
	var params utils.HTTPRequestParams
	params.URL = "https://github-readme-stats.vercel.app/api?username=" + github
	params.Method = "GET"
	content, err := utils.SendHTTPRequest(params)
	if err != nil {
		logger.Log.Errorf("Request github stat page errror, %s, %s", params.URL, err.Error())
		return "", errors.New("Can't get GitHub Rank")
	}

	rank, err := utils.GetGitRank(content)
	if err != nil {
		logger.Log.Errorf("Parse GitHub's rank error, %s, %s", params.URL, err.Error())
		return "", errors.New("Can't parse GitHub Rank")
	}

	var amount string
	switch rank {
	case "S":
		amount = "1"
	case "A":
		amount = "0.4"
	case "B":
		amount = "0.3"
	case "C":
		amount = "0.1"
	default:
		logger.Log.Errorf("github's rank is invalid, %s, %s", params.URL, rank)
		return "", errors.New("GitHub's rank is invalid")
	}
	return amount, nil
}

func RequestGitRankWithToken(github string, token string) (string, error) {
	var params utils.HTTPRequestParams
	params.URL = "https://github-readme-stats.vercel.app/api?username=" + github
	params.Method = "GET"
	content, err := utils.SendHTTPRequest(params)
	if err != nil {
		logger.Log.Errorf("Request github stat page errror, %s, %s", params.URL, err.Error())
		return "", errors.New("Can't get GitHub Rank")
	}

	rank, err := utils.GetGitRank(content)
	if err != nil {
		logger.Log.Errorf("Parse GitHub's rank error, %s, %s", params.URL, err.Error())
		return "", errors.New("Can't parse GitHub Rank")
	}

	var amount string
	switch rank {
	case "S":
		amount = "5"
	case "A":
		amount = "4"
	case "B":
		amount = "3.5"
	case "C":
		amount = "3"
	default:
		logger.Log.Errorf("github's rank is invalid, %s, %s", params.URL, rank)
		return "", errors.New("GitHub's rank is invalid")
	}
	return amount, nil
}

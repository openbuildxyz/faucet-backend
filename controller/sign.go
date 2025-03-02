package controller

import (
	"encoding/json"
	"faucet/logger"
	"faucet/model"
	"faucet/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func HandleSign(c *gin.Context) {
	var req SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Errorf("Invalid request: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request. Please try again later.", nil)
		return
	}

	var accessRequest AccessTokenRequest
	accessRequest.ClientId = viper.GetString("oauth.clientId")
	accessRequest.ClientSecret = viper.GetString("oauth.clientSecret")
	accessRequest.Code = req.Code

	var reqArgs utils.HTTPRequestParams
	reqArgs.URL = viper.GetString("oauth.accessApi")
	reqArgs.Method = "POST"
	reqArgs.Body = accessRequest

	result, err := utils.SendHTTPRequest(reqArgs)
	if err != nil {
		logger.Log.Errorf("ServerError: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	var resp AccessTokenResponse
	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		logger.Log.Errorf("ServerError: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	if resp.Status != 200 {
		logger.Log.Errorf("ServerError: %v", resp)
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	var user model.User
	// oauth
	user.OauthToken = resp.Data.Token
	err = model.CreateUser(&user)
	if err != nil {
		logger.Log.Errorf("ServerError: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	// gin token
	token, err := utils.GenerateToken(resp.Data.Token)
	if err != nil {
		logger.Log.Errorf("GenerateToken err: %s", err.Error())
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	var response SignResponse
	// oauth token
	response.Token = token

	utils.SuccessResponse(c, http.StatusOK, "success", response)
}

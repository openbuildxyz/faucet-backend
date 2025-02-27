package controller

import (
	"encoding/json"
	"faucet/logger"
	"faucet/model"
	"faucet/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func HandleGetUser(c *gin.Context) {
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

	var reqArgs utils.HTTPRequestParams
	reqArgs.URL = viper.GetString("oauth.getUser")
	reqArgs.Method = "GET"

	header := make(map[string]string)
	header["Authorization"] = authHeader
	reqArgs.Headers = header

	result, err := utils.SendHTTPRequest(reqArgs)
	if err != nil {
		logger.Log.Errorf("ServerError: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var resp GetUserResponse
	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		logger.Log.Errorf("ServerError: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if resp.Status != 200 {
		logger.Log.Errorf("ServerError: %v", resp)
		utils.ErrorResponse(c, http.StatusInternalServerError, resp.Message, nil)
		return
	}

	fmt.Println(resp)

	user.Uid = resp.Data.Uid
	user.Avatar = resp.Data.Avatar
	user.Email = resp.Data.Email
	user.Username = resp.Data.UserName
	user.Github = resp.Data.Github
	err = model.UpdateUser(user)
	if err != nil {
		logger.Log.Errorf("ServerError: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", user)
}

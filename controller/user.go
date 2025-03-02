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
		utils.ErrorResponse(c, http.StatusUnauthorized, "", nil)
		return
	}

	var reqArgs utils.HTTPRequestParams
	reqArgs.URL = viper.GetString("oauth.getUser")
	reqArgs.Method = "GET"

	header := make(map[string]string)
	header["Authorization"] = fmt.Sprintf("Bearer %s", oToken)
	reqArgs.Headers = header

	result, err := utils.SendHTTPRequest(reqArgs)
	if err != nil {
		logger.Log.Errorf("SendHTTPRequest err: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	var resp GetUserResponse
	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		logger.Log.Errorf("Unmarshal err: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	if resp.Status != 200 {
		logger.Log.Errorf("ServerError: %v", resp)
		utils.ErrorResponse(c, http.StatusInternalServerError, resp.Message, nil)
		return
	}

	user.Uid = resp.Data.Uid
	user.Avatar = resp.Data.Avatar
	user.Email = resp.Data.Email
	user.Username = resp.Data.UserName
	user.Github = resp.Data.Github
	user.OauthTokenId = resp.ID
	err = model.UpdateUser(user)
	if err != nil {
		logger.Log.Errorf("ServerError: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "The system is currently busy. Please try again later.", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", user)
}

package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func FaucetHandler(c *gin.Context) {
}

func sendEther(client *ethclient.Client, privateKeyHex string, to common.Address, amount string, c *gin.Context) error {
	return nil
}

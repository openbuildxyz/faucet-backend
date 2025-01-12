package model

import (
	"faucet/config"
)

var db = config.DB

func init() {
	db.AutoMigrate(&Transaction{})
}

package chain

import (
	"faucet/model"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

func Transfer(receive string, value *big.Int) (string, error) {
	var trans model.Transaction

	_, MaxPriority, MaxFee, err := GetGasPrice()
	trans.MaxPriorityFeePerGas = MaxPriority
	trans.MaxFeePerGas = MaxFee
	to := common.HexToAddress(receive)
	trans.ContractAddress = &to
	trans.Value = value

	sk := viper.GetString("faucet.private_key")
	trans.SK = sk

	trans.GasLimit = 21000

	tx, err := SendTransaction(trans)
	if err != nil {
		return "", err
	}

	go CheckTransaction(tx)

	fmt.Println("提交转账交易: ", tx)
	return tx, nil
}

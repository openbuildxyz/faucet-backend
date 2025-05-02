package chain

import (
	"crypto/ecdsa"
	"errors"
	"faucet/model"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

func Transfer(receive string, token string, value *big.Int) (string, error) {
	var trans model.RawTransaction

	_, _, MaxFee, err := GetGasPrice()
	trans.MaxPriorityFeePerGas = big.NewInt(52)
	trans.MaxFeePerGas = MaxFee
	to := common.HexToAddress(receive)
	trans.To = &to
	trans.Value = value

	sk := viper.GetString("faucet.private_key")
	trans.SK = sk

	privateKey, err := crypto.HexToECDSA(sk)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("error casting public key to ECDSA")
	}
	from := crypto.PubkeyToAddress(*publicKeyECDSA)

	gasLimit, err := GetGasLimit(&from, &to, trans.Data, trans.Value)
	if err != nil {
		return "", err
	}
	trans.GasLimit = gasLimit

	var chainInfo string
	if token == "MON" {
		chainInfo = "rpc.MonadDevnet"
	}
	if token == "0G" {
		chainInfo = "rpc.ZeroTestnet"
	}

	ReconnetRpc(chainInfo)

	tx, err := SendTransaction(trans)
	return tx, err
}

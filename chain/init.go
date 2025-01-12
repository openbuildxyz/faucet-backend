package chain

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
)

var EthClient *ethclient.Client
var RpcClient *rpc.Client

func init() {
	rpcUrl := viper.GetString("rpc.EthSepolia")

	fmt.Println(rpcUrl)

	ethClient, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	EthClient = ethClient

	// rpc
	rpcClient, err := rpc.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum rpc: %v", err)
	}

	RpcClient = rpcClient
}

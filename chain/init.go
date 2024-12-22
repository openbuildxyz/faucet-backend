package chain

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var EthClient *ethclient.Client
var RpcClient *rpc.Client

var rpcUrl = "https://eth-sepolia.g.alchemy.com/v2/YQHOAbPqR8tHixcwD-ZW5xxUtVjbEmHA"

// var rpcUrl = "https://mainnet.infura.io/v3/672e64bfa6f144349608236513a79679"

func init() {
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

package chain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"faucet/model"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func SendTransaction(transaction model.Transaction) (string, error) {
	privateKey, err := crypto.HexToECDSA(transaction.SK)
	if err != nil {
		// log.Fatal(err)
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		// log.Fatal("error casting public key to ECDSA")
		return "", errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := EthClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		// log.Fatal(err)
		return "", err
	}

	chainID, err := EthClient.NetworkID(context.Background())

	// 使用types.NewTx创建EIP-1559类型的交易
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: transaction.MaxPriorityFeePerGas, // 设置合适的 MaxPriorityFeePerGas
		GasFeeCap: transaction.MaxFeePerGas,         // 设置合适的 MaxFeePerGas
		Gas:       transaction.GasLimit,             // 设置合适的 Gas Limit
		To:        transaction.ContractAddress,
		Value:     transaction.Value,
		Data:      transaction.Data,
	})

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		// log.Fatal(err)
		return "", err
	}

	err = EthClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		// log.Fatal(err)
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// TODO:
// MQ提高稳定性
// DB数据持久化
func CheckTransaction(tx string) (bool, error) {
	txHash := common.HexToHash(tx)
	time.Sleep(10 * time.Second)

	var checkcount = 0

	// 轮询检查交易是否被确认
	for {
		time.Sleep(5 * time.Second)
		checkcount++
		if checkcount >= 20 {
			return false, nil
		}
		time.Sleep(10 * time.Second)
		_, isPending, err := EthClient.TransactionByHash(context.Background(), txHash)
		if err != nil {
			fmt.Println("Error getting transaction: ", err)
			errInfo := fmt.Sprintf("%s", err.Error())
			if strings.Contains(errInfo, "not found") {
				continue
			}
			fmt.Println("CheckTransaction TransactionByHash err: ", err)
			return false, err
		}

		if !isPending {
			receipt, err := EthClient.TransactionReceipt(context.Background(), txHash)
			if err != nil {
				errInfo := fmt.Sprintf("%s", err.Error())
				if strings.Contains(errInfo, "not found") {
					continue
				}
				// log.Fatalf("Error getting transaction receipt: %v", err)

				fmt.Println("CheckTransaction TransactionReceipt err: ", err)
				return false, err
			}

			if receipt.Status == types.ReceiptStatusSuccessful {
				// fmt.Println("Transaction confirmed!")
				return true, nil
			} else {
				// fmt.Println("Transaction failed!")
				return false, nil
			}
		}

	}
}

func GetGasPrice() (*big.Int, *big.Int, *big.Int, error) {
	var block map[string]interface{}
	err := RpcClient.CallContext(context.Background(), &block, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		// panic(err)
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
	}
	baseFeePerGasStr := block["baseFeePerGas"].(string) // 获取baseFeePerGas为字符串
	baseFeePerGas := new(big.Int)
	baseFeePerGas, ok := baseFeePerGas.SetString(baseFeePerGasStr[2:], 16) // 从十六进制转换为big.Int
	if !ok {
		// panic("Failed to parse baseFeePerGas")
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), errors.New("Failed to parse baseFeePerGas")
	}

	// 设置最大优先费用（这里示例设置为2 Gwei）
	maxPriorityFeePerGas := big.NewInt(2e9) // 2 Gwei
	// fmt.Printf("Max Priority Fee Per Gas: %s\n", maxPriorityFeePerGas.String())

	// 设置最大费用，假设为基本费用的2倍加上最大优先费用
	maxFeePerGas := new(big.Int).Mul(baseFeePerGas, big.NewInt(3))
	maxFeePerGas = maxFeePerGas.Add(maxFeePerGas, maxPriorityFeePerGas)

	return baseFeePerGas, maxPriorityFeePerGas, maxFeePerGas, nil
}

func GetGasLimit(fromAddress, toAddress *common.Address, data []byte, value *big.Int) uint64 {
	// 构造交易数据
	tx := map[string]interface{}{
		"from":  fromAddress,
		"to":    toAddress,
		"data":  "0x" + hex.EncodeToString(data),
		"value": hexutil.EncodeBig(value),
	}

	var gasLimit hexutil.Uint64
	if err := RpcClient.CallContext(context.Background(), &gasLimit, "eth_estimateGas", tx); err != nil {
		fmt.Printf("Failed to estimate gas: %v", err)
	}

	fmt.Printf("Estimated Gas: %d\n", gasLimit)
	return uint64(gasLimit)
}

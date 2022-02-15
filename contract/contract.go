package contract

import (
	"io/ioutil"
	"math/big"

	web3 "github.com/umbracle/go-web3"
	web3abi "github.com/umbracle/go-web3/abi"
	web3contract "github.com/umbracle/go-web3/contract"
	"github.com/umbracle/go-web3/jsonrpc"
)

func GetContract(
	provider *jsonrpc.Client,
	contractAddress web3.Address,
	abiPath string,
) (*web3contract.Contract, error) {
	abi, err := loadAbi(abiPath)
	if err != nil {
		return nil, err
	}
	return web3contract.NewContract(contractAddress, abi, provider), nil
}

func DecodeBigInt(v map[string]interface{}, key string) (float64, bool) {
	val1, ok := v[key]
	if !ok {
		return 0, false
	}

	val2, ok := val1.(*big.Int)
	if !ok {
		return 0, false
	}

	fnumber, _ := new(big.Float).SetInt(val2).Float64()
	return fnumber, true
}

func loadAbi(jsonPath string) (*web3abi.ABI, error) {
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	return web3abi.NewABI(string(data))
}

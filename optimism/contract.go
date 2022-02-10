package optimism

import (
	"io/ioutil"
	"math/big"
	"path/filepath"

	web3 "github.com/umbracle/go-web3"
	web3abi "github.com/umbracle/go-web3/abi"
	web3contract "github.com/umbracle/go-web3/contract"
	"github.com/umbracle/go-web3/jsonrpc"
)

type OptimismClient struct {
	provider *jsonrpc.Client
	abiPath  string
}

func NewOptimismClient(provider *jsonrpc.Client, abiPath string) *OptimismClient {
	return &OptimismClient{provider, abiPath}
}

func (p *OptimismClient) GetContract(contractName string, contractAddress web3.Address) (*web3contract.Contract, error) {
	abi, err := loadAbi(filepath.Join(p.abiPath, contractName+".json"))
	if err != nil {
		return nil, err
	}
	return web3contract.NewContract(contractAddress, abi, p.provider), nil
}

func (p *OptimismClient) DecodeBigInt(v map[string]interface{}, key string) (float64, bool) {
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

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multicall2

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Multicall2Call is an auto generated low-level Go binding around an user-defined struct.
type Multicall2Call struct {
	Target   common.Address
	CallData []byte
}

// Multicall2Result is an auto generated low-level Go binding around an user-defined struct.
type Multicall2Result struct {
	Success    bool
	ReturnData []byte
}

// Multicall2MetaData contains all meta data concerning the Multicall2 contract.
var Multicall2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall2.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"aggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"returnData\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall2.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"blockAndAggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall2.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"getBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockCoinbase\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"coinbase\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockDifficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"difficulty\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockGasLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"gaslimit\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"requireSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall2.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"tryAggregate\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall2.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"requireSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall2.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"tryBlockAndAggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall2.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610f02806100206000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c806372425d9d1161007157806372425d9d146101a657806386d516e8146101c4578063a8b0574e146101e2578063bce38bd714610200578063c3077fa914610230578063ee82ac5e14610262576100b4565b80630f28c97d146100b9578063252dba42146100d757806327e86d6e14610108578063399542e91461012657806342cbb15c146101585780634d2301cc14610176575b600080fd5b6100c1610292565b6040516100ce9190610523565b60405180910390f35b6100f160048036038101906100ec9190610852565b61029a565b6040516100ff9291906109e5565b60405180910390f35b6101106102a3565b60405161011d9190610a2e565b60405180910390f35b610140600480360381019061013b9190610a81565b6102b8565b60405161014f93929190610beb565b60405180910390f35b6101606102d7565b60405161016d9190610523565b60405180910390f35b610190600480360381019061018b9190610c29565b6102df565b60405161019d9190610523565b60405180910390f35b6101ae610300565b6040516101bb9190610523565b60405180910390f35b6101cc610308565b6040516101d99190610523565b60405180910390f35b6101ea610310565b6040516101f79190610c65565b60405180910390f35b61021a60048036038101906102159190610a81565b610318565b6040516102279190610c80565b60405180910390f35b61024a60048036038101906102459190610852565b6104c0565b60405161025993929190610beb565b60405180910390f35b61027c60048036038101906102779190610cce565b6104e3565b6040516102899190610a2e565b60405180910390f35b600042905090565b60006060915091565b60006001436102b29190610d2a565b40905090565b6000806060439250434091506102ce8585610318565b90509250925092565b600043905090565b60008173ffffffffffffffffffffffffffffffffffffffff16319050919050565b600044905090565b600045905090565b600041905090565b6060815167ffffffffffffffff81111561033557610334610568565b5b60405190808252806020026020018201604052801561036e57816020015b61035b6104ee565b8152602001906001900390816103535790505b50905060005b82518110156104b95760008084838151811061039357610392610d5e565b5b60200260200101516000015173ffffffffffffffffffffffffffffffffffffffff168584815181106103c8576103c7610d5e565b5b6020026020010151602001516040516103e19190610dc9565b6000604051808303816000865af19150503d806000811461041e576040519150601f19603f3d011682016040523d82523d6000602084013e610423565b606091505b5091509150851561046f578161046e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161046590610e63565b60405180910390fd5b5b604051806040016040528083151581526020018281525084848151811061049957610498610d5e565b5b6020026020010181905250505080806104b190610e83565b915050610374565b5092915050565b60008060606104d06001856102b8565b8093508194508295505050509193909250565b600081409050919050565b6040518060400160405280600015158152602001606081525090565b6000819050919050565b61051d8161050a565b82525050565b60006020820190506105386000830184610514565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105a082610557565b810181811067ffffffffffffffff821117156105bf576105be610568565b5b80604052505050565b60006105d261053e565b90506105de8282610597565b919050565b600067ffffffffffffffff8211156105fe576105fd610568565b5b602082029050602081019050919050565b600080fd5b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006106498261061e565b9050919050565b6106598161063e565b811461066457600080fd5b50565b60008135905061067681610650565b92915050565b600080fd5b600067ffffffffffffffff82111561069c5761069b610568565b5b6106a582610557565b9050602081019050919050565b82818337600083830152505050565b60006106d46106cf84610681565b6105c8565b9050828152602081018484840111156106f0576106ef61067c565b5b6106fb8482856106b2565b509392505050565b600082601f83011261071857610717610552565b5b81356107288482602086016106c1565b91505092915050565b60006040828403121561074757610746610614565b5b61075160406105c8565b9050600061076184828501610667565b600083015250602082013567ffffffffffffffff81111561078557610784610619565b5b61079184828501610703565b60208301525092915050565b60006107b06107ab846105e3565b6105c8565b905080838252602082019050602084028301858111156107d3576107d261060f565b5b835b8181101561081a57803567ffffffffffffffff8111156107f8576107f7610552565b5b8086016108058982610731565b855260208501945050506020810190506107d5565b5050509392505050565b600082601f83011261083957610838610552565b5b813561084984826020860161079d565b91505092915050565b60006020828403121561086857610867610548565b5b600082013567ffffffffffffffff8111156108865761088561054d565b5b61089284828501610824565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b838110156109015780820151818401526020810190506108e6565b83811115610910576000848401525b50505050565b6000610921826108c7565b61092b81856108d2565b935061093b8185602086016108e3565b61094481610557565b840191505092915050565b600061095b8383610916565b905092915050565b6000602082019050919050565b600061097b8261089b565b61098581856108a6565b935083602082028501610997856108b7565b8060005b858110156109d357848403895281516109b4858261094f565b94506109bf83610963565b925060208a0199505060018101905061099b565b50829750879550505050505092915050565b60006040820190506109fa6000830185610514565b8181036020830152610a0c8184610970565b90509392505050565b6000819050919050565b610a2881610a15565b82525050565b6000602082019050610a436000830184610a1f565b92915050565b60008115159050919050565b610a5e81610a49565b8114610a6957600080fd5b50565b600081359050610a7b81610a55565b92915050565b60008060408385031215610a9857610a97610548565b5b6000610aa685828601610a6c565b925050602083013567ffffffffffffffff811115610ac757610ac661054d565b5b610ad385828601610824565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b610b1281610a49565b82525050565b6000604083016000830151610b306000860182610b09565b5060208301518482036020860152610b488282610916565b9150508091505092915050565b6000610b618383610b18565b905092915050565b6000602082019050919050565b6000610b8182610add565b610b8b8185610ae8565b935083602082028501610b9d85610af9565b8060005b85811015610bd95784840389528151610bba8582610b55565b9450610bc583610b69565b925060208a01995050600181019050610ba1565b50829750879550505050505092915050565b6000606082019050610c006000830186610514565b610c0d6020830185610a1f565b8181036040830152610c1f8184610b76565b9050949350505050565b600060208284031215610c3f57610c3e610548565b5b6000610c4d84828501610667565b91505092915050565b610c5f8161063e565b82525050565b6000602082019050610c7a6000830184610c56565b92915050565b60006020820190508181036000830152610c9a8184610b76565b905092915050565b610cab8161050a565b8114610cb657600080fd5b50565b600081359050610cc881610ca2565b92915050565b600060208284031215610ce457610ce3610548565b5b6000610cf284828501610cb9565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610d358261050a565b9150610d408361050a565b925082821015610d5357610d52610cfb565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081905092915050565b6000610da3826108c7565b610dad8185610d8d565b9350610dbd8185602086016108e3565b80840191505092915050565b6000610dd58284610d98565b915081905092915050565b600082825260208201905092915050565b7f4d756c746963616c6c32206167677265676174653a2063616c6c206661696c6560008201527f6400000000000000000000000000000000000000000000000000000000000000602082015250565b6000610e4d602183610de0565b9150610e5882610df1565b604082019050919050565b60006020820190508181036000830152610e7c81610e40565b9050919050565b6000610e8e8261050a565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610ec157610ec0610cfb565b5b60018201905091905056fea2646970667358221220eda033ddd87162c0eba63801d88ceb44e5e13e40a118d0bea2e7d867e562398a64736f6c63430008090033",
}

// Multicall2ABI is the input ABI used to generate the binding from.
// Deprecated: Use Multicall2MetaData.ABI instead.
var Multicall2ABI = Multicall2MetaData.ABI

// Multicall2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Multicall2MetaData.Bin instead.
var Multicall2Bin = Multicall2MetaData.Bin

// DeployMulticall2 deploys a new Ethereum contract, binding an instance of Multicall2 to it.
func DeployMulticall2(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Multicall2, error) {
	parsed, err := Multicall2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Multicall2Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Multicall2{Multicall2Caller: Multicall2Caller{contract: contract}, Multicall2Transactor: Multicall2Transactor{contract: contract}, Multicall2Filterer: Multicall2Filterer{contract: contract}}, nil
}

// Multicall2 is an auto generated Go binding around an Ethereum contract.
type Multicall2 struct {
	Multicall2Caller     // Read-only binding to the contract
	Multicall2Transactor // Write-only binding to the contract
	Multicall2Filterer   // Log filterer for contract events
}

// Multicall2Caller is an auto generated read-only Go binding around an Ethereum contract.
type Multicall2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Multicall2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Multicall2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Multicall2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Multicall2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Multicall2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Multicall2Session struct {
	Contract     *Multicall2       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Multicall2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Multicall2CallerSession struct {
	Contract *Multicall2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// Multicall2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Multicall2TransactorSession struct {
	Contract     *Multicall2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// Multicall2Raw is an auto generated low-level Go binding around an Ethereum contract.
type Multicall2Raw struct {
	Contract *Multicall2 // Generic contract binding to access the raw methods on
}

// Multicall2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Multicall2CallerRaw struct {
	Contract *Multicall2Caller // Generic read-only contract binding to access the raw methods on
}

// Multicall2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Multicall2TransactorRaw struct {
	Contract *Multicall2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewMulticall2 creates a new instance of Multicall2, bound to a specific deployed contract.
func NewMulticall2(address common.Address, backend bind.ContractBackend) (*Multicall2, error) {
	contract, err := bindMulticall2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Multicall2{Multicall2Caller: Multicall2Caller{contract: contract}, Multicall2Transactor: Multicall2Transactor{contract: contract}, Multicall2Filterer: Multicall2Filterer{contract: contract}}, nil
}

// NewMulticall2Caller creates a new read-only instance of Multicall2, bound to a specific deployed contract.
func NewMulticall2Caller(address common.Address, caller bind.ContractCaller) (*Multicall2Caller, error) {
	contract, err := bindMulticall2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Multicall2Caller{contract: contract}, nil
}

// NewMulticall2Transactor creates a new write-only instance of Multicall2, bound to a specific deployed contract.
func NewMulticall2Transactor(address common.Address, transactor bind.ContractTransactor) (*Multicall2Transactor, error) {
	contract, err := bindMulticall2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Multicall2Transactor{contract: contract}, nil
}

// NewMulticall2Filterer creates a new log filterer instance of Multicall2, bound to a specific deployed contract.
func NewMulticall2Filterer(address common.Address, filterer bind.ContractFilterer) (*Multicall2Filterer, error) {
	contract, err := bindMulticall2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Multicall2Filterer{contract: contract}, nil
}

// bindMulticall2 binds a generic wrapper to an already deployed contract.
func bindMulticall2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Multicall2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multicall2 *Multicall2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multicall2.Contract.Multicall2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multicall2 *Multicall2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multicall2.Contract.Multicall2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multicall2 *Multicall2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multicall2.Contract.Multicall2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multicall2 *Multicall2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multicall2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multicall2 *Multicall2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multicall2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multicall2 *Multicall2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multicall2.Contract.contract.Transact(opts, method, params...)
}

// Aggregate is a free data retrieval call binding the contract method 0x252dba42.
//
// Solidity: function aggregate((address,bytes)[] calls) view returns(uint256 blockNumber, bytes[] returnData)
func (_Multicall2 *Multicall2Caller) Aggregate(opts *bind.CallOpts, calls []Multicall2Call) (struct {
	BlockNumber *big.Int
	ReturnData  [][]byte
}, error) {
	var out []interface{}
	err := _Multicall2.contract.Call(opts, &out, "aggregate", calls)

	outstruct := new(struct {
		BlockNumber *big.Int
		ReturnData  [][]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockNumber = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ReturnData = *abi.ConvertType(out[1], new([][]byte)).(*[][]byte)

	return *outstruct, err

}

// Aggregate is a free data retrieval call binding the contract method 0x252dba42.
//
// Solidity: function aggregate((address,bytes)[] calls) view returns(uint256 blockNumber, bytes[] returnData)
func (_Multicall2 *Multicall2Session) Aggregate(calls []Multicall2Call) (struct {
	BlockNumber *big.Int
	ReturnData  [][]byte
}, error) {
	return _Multicall2.Contract.Aggregate(&_Multicall2.CallOpts, calls)
}

// Aggregate is a free data retrieval call binding the contract method 0x252dba42.
//
// Solidity: function aggregate((address,bytes)[] calls) view returns(uint256 blockNumber, bytes[] returnData)
func (_Multicall2 *Multicall2CallerSession) Aggregate(calls []Multicall2Call) (struct {
	BlockNumber *big.Int
	ReturnData  [][]byte
}, error) {
	return _Multicall2.Contract.Aggregate(&_Multicall2.CallOpts, calls)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 blockNumber) view returns(bytes32 blockHash)
func (_Multicall2 *Multicall2Caller) GetBlockHash(opts *bind.CallOpts, blockNumber *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Multicall2.contract.Call(opts, &out, "getBlockHash", blockNumber)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBlockHash is a free data retrieval call binding the contract method 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 blockNumber) view returns(bytes32 blockHash)
func (_Multicall2 *Multicall2Session) GetBlockHash(blockNumber *big.Int) ([32]byte, error) {
	return _Multicall2.Contract.GetBlockHash(&_Multicall2.CallOpts, blockNumber)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 blockNumber) view returns(bytes32 blockHash)
func (_Multicall2 *Multicall2CallerSession) GetBlockHash(blockNumber *big.Int) ([32]byte, error) {
	return _Multicall2.Contract.GetBlockHash(&_Multicall2.CallOpts, blockNumber)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint256 blockNumber)
func (_Multicall2 *Multicall2Caller) GetBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall2.contract.Call(opts, &out, "getBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint256 blockNumber)
func (_Multicall2 *Multicall2Session) GetBlockNumber() (*big.Int, error) {
	return _Multicall2.Contract.GetBlockNumber(&_Multicall2.CallOpts)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint256 blockNumber)
func (_Multicall2 *Multicall2CallerSession) GetBlockNumber() (*big.Int, error) {
	return _Multicall2.Contract.GetBlockNumber(&_Multicall2.CallOpts)
}

// GetCurrentBlockCoinbase is a free data retrieval call binding the contract method 0xa8b0574e.
//
// Solidity: function getCurrentBlockCoinbase() view returns(address coinbase)
func (_Multicall2 *Multicall2Caller) GetCurrentBlockCoinbase(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Multicall2.contract.Call(opts, &out, "getCurrentBlockCoinbase")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetCurrentBlockCoinbase is a free data retrieval call binding the contract method 0xa8b0574e.
//
// Solidity: function getCurrentBlockCoinbase() view returns(address coinbase)
func (_Multicall2 *Multicall2Session) GetCurrentBlockCoinbase() (common.Address, error) {
	return _Multicall2.Contract.GetCurrentBlockCoinbase(&_Multicall2.CallOpts)
}

// GetCurrentBlockCoinbase is a free data retrieval call binding the contract method 0xa8b0574e.
//
// Solidity: function getCurrentBlockCoinbase() view returns(address coinbase)
func (_Multicall2 *Multicall2CallerSession) GetCurrentBlockCoinbase() (common.Address, error) {
	return _Multicall2.Contract.GetCurrentBlockCoinbase(&_Multicall2.CallOpts)
}

// GetCurrentBlockDifficulty is a free data retrieval call binding the contract method 0x72425d9d.
//
// Solidity: function getCurrentBlockDifficulty() view returns(uint256 difficulty)
func (_Multicall2 *Multicall2Caller) GetCurrentBlockDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall2.contract.Call(opts, &out, "getCurrentBlockDifficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentBlockDifficulty is a free data retrieval call binding the contract method 0x72425d9d.
//
// Solidity: function getCurrentBlockDifficulty() view returns(uint256 difficulty)
func (_Multicall2 *Multicall2Session) GetCurrentBlockDifficulty() (*big.Int, error) {
	return _Multicall2.Contract.GetCurrentBlockDifficulty(&_Multicall2.CallOpts)
}

// GetCurrentBlockDifficulty is a free data retrieval call binding the contract method 0x72425d9d.
//
// Solidity: function getCurrentBlockDifficulty() view returns(uint256 difficulty)
func (_Multicall2 *Multicall2CallerSession) GetCurrentBlockDifficulty() (*big.Int, error) {
	return _Multicall2.Contract.GetCurrentBlockDifficulty(&_Multicall2.CallOpts)
}

// GetCurrentBlockGasLimit is a free data retrieval call binding the contract method 0x86d516e8.
//
// Solidity: function getCurrentBlockGasLimit() view returns(uint256 gaslimit)
func (_Multicall2 *Multicall2Caller) GetCurrentBlockGasLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall2.contract.Call(opts, &out, "getCurrentBlockGasLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentBlockGasLimit is a free data retrieval call binding the contract method 0x86d516e8.
//
// Solidity: function getCurrentBlockGasLimit() view returns(uint256 gaslimit)
func (_Multicall2 *Multicall2Session) GetCurrentBlockGasLimit() (*big.Int, error) {
	return _Multicall2.Contract.GetCurrentBlockGasLimit(&_Multicall2.CallOpts)
}

// GetCurrentBlockGasLimit is a free data retrieval call binding the contract method 0x86d516e8.
//
// Solidity: function getCurrentBlockGasLimit() view returns(uint256 gaslimit)
func (_Multicall2 *Multicall2CallerSession) GetCurrentBlockGasLimit() (*big.Int, error) {
	return _Multicall2.Contract.GetCurrentBlockGasLimit(&_Multicall2.CallOpts)
}

// GetCurrentBlockTimestamp is a free data retrieval call binding the contract method 0x0f28c97d.
//
// Solidity: function getCurrentBlockTimestamp() view returns(uint256 timestamp)
func (_Multicall2 *Multicall2Caller) GetCurrentBlockTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multicall2.contract.Call(opts, &out, "getCurrentBlockTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentBlockTimestamp is a free data retrieval call binding the contract method 0x0f28c97d.
//
// Solidity: function getCurrentBlockTimestamp() view returns(uint256 timestamp)
func (_Multicall2 *Multicall2Session) GetCurrentBlockTimestamp() (*big.Int, error) {
	return _Multicall2.Contract.GetCurrentBlockTimestamp(&_Multicall2.CallOpts)
}

// GetCurrentBlockTimestamp is a free data retrieval call binding the contract method 0x0f28c97d.
//
// Solidity: function getCurrentBlockTimestamp() view returns(uint256 timestamp)
func (_Multicall2 *Multicall2CallerSession) GetCurrentBlockTimestamp() (*big.Int, error) {
	return _Multicall2.Contract.GetCurrentBlockTimestamp(&_Multicall2.CallOpts)
}

// GetEthBalance is a free data retrieval call binding the contract method 0x4d2301cc.
//
// Solidity: function getEthBalance(address addr) view returns(uint256 balance)
func (_Multicall2 *Multicall2Caller) GetEthBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Multicall2.contract.Call(opts, &out, "getEthBalance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEthBalance is a free data retrieval call binding the contract method 0x4d2301cc.
//
// Solidity: function getEthBalance(address addr) view returns(uint256 balance)
func (_Multicall2 *Multicall2Session) GetEthBalance(addr common.Address) (*big.Int, error) {
	return _Multicall2.Contract.GetEthBalance(&_Multicall2.CallOpts, addr)
}

// GetEthBalance is a free data retrieval call binding the contract method 0x4d2301cc.
//
// Solidity: function getEthBalance(address addr) view returns(uint256 balance)
func (_Multicall2 *Multicall2CallerSession) GetEthBalance(addr common.Address) (*big.Int, error) {
	return _Multicall2.Contract.GetEthBalance(&_Multicall2.CallOpts, addr)
}

// GetLastBlockHash is a free data retrieval call binding the contract method 0x27e86d6e.
//
// Solidity: function getLastBlockHash() view returns(bytes32 blockHash)
func (_Multicall2 *Multicall2Caller) GetLastBlockHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Multicall2.contract.Call(opts, &out, "getLastBlockHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetLastBlockHash is a free data retrieval call binding the contract method 0x27e86d6e.
//
// Solidity: function getLastBlockHash() view returns(bytes32 blockHash)
func (_Multicall2 *Multicall2Session) GetLastBlockHash() ([32]byte, error) {
	return _Multicall2.Contract.GetLastBlockHash(&_Multicall2.CallOpts)
}

// GetLastBlockHash is a free data retrieval call binding the contract method 0x27e86d6e.
//
// Solidity: function getLastBlockHash() view returns(bytes32 blockHash)
func (_Multicall2 *Multicall2CallerSession) GetLastBlockHash() ([32]byte, error) {
	return _Multicall2.Contract.GetLastBlockHash(&_Multicall2.CallOpts)
}

// BlockAndAggregate is a paid mutator transaction binding the contract method 0xc3077fa9.
//
// Solidity: function blockAndAggregate((address,bytes)[] calls) returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (_Multicall2 *Multicall2Transactor) BlockAndAggregate(opts *bind.TransactOpts, calls []Multicall2Call) (*types.Transaction, error) {
	return _Multicall2.contract.Transact(opts, "blockAndAggregate", calls)
}

// BlockAndAggregate is a paid mutator transaction binding the contract method 0xc3077fa9.
//
// Solidity: function blockAndAggregate((address,bytes)[] calls) returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (_Multicall2 *Multicall2Session) BlockAndAggregate(calls []Multicall2Call) (*types.Transaction, error) {
	return _Multicall2.Contract.BlockAndAggregate(&_Multicall2.TransactOpts, calls)
}

// BlockAndAggregate is a paid mutator transaction binding the contract method 0xc3077fa9.
//
// Solidity: function blockAndAggregate((address,bytes)[] calls) returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (_Multicall2 *Multicall2TransactorSession) BlockAndAggregate(calls []Multicall2Call) (*types.Transaction, error) {
	return _Multicall2.Contract.BlockAndAggregate(&_Multicall2.TransactOpts, calls)
}

// TryAggregate is a paid mutator transaction binding the contract method 0xbce38bd7.
//
// Solidity: function tryAggregate(bool requireSuccess, (address,bytes)[] calls) returns((bool,bytes)[] returnData)
func (_Multicall2 *Multicall2Transactor) TryAggregate(opts *bind.TransactOpts, requireSuccess bool, calls []Multicall2Call) (*types.Transaction, error) {
	return _Multicall2.contract.Transact(opts, "tryAggregate", requireSuccess, calls)
}

// TryAggregate is a paid mutator transaction binding the contract method 0xbce38bd7.
//
// Solidity: function tryAggregate(bool requireSuccess, (address,bytes)[] calls) returns((bool,bytes)[] returnData)
func (_Multicall2 *Multicall2Session) TryAggregate(requireSuccess bool, calls []Multicall2Call) (*types.Transaction, error) {
	return _Multicall2.Contract.TryAggregate(&_Multicall2.TransactOpts, requireSuccess, calls)
}

// TryAggregate is a paid mutator transaction binding the contract method 0xbce38bd7.
//
// Solidity: function tryAggregate(bool requireSuccess, (address,bytes)[] calls) returns((bool,bytes)[] returnData)
func (_Multicall2 *Multicall2TransactorSession) TryAggregate(requireSuccess bool, calls []Multicall2Call) (*types.Transaction, error) {
	return _Multicall2.Contract.TryAggregate(&_Multicall2.TransactOpts, requireSuccess, calls)
}

// TryBlockAndAggregate is a paid mutator transaction binding the contract method 0x399542e9.
//
// Solidity: function tryBlockAndAggregate(bool requireSuccess, (address,bytes)[] calls) returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (_Multicall2 *Multicall2Transactor) TryBlockAndAggregate(opts *bind.TransactOpts, requireSuccess bool, calls []Multicall2Call) (*types.Transaction, error) {
	return _Multicall2.contract.Transact(opts, "tryBlockAndAggregate", requireSuccess, calls)
}

// TryBlockAndAggregate is a paid mutator transaction binding the contract method 0x399542e9.
//
// Solidity: function tryBlockAndAggregate(bool requireSuccess, (address,bytes)[] calls) returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (_Multicall2 *Multicall2Session) TryBlockAndAggregate(requireSuccess bool, calls []Multicall2Call) (*types.Transaction, error) {
	return _Multicall2.Contract.TryBlockAndAggregate(&_Multicall2.TransactOpts, requireSuccess, calls)
}

// TryBlockAndAggregate is a paid mutator transaction binding the contract method 0x399542e9.
//
// Solidity: function tryBlockAndAggregate(bool requireSuccess, (address,bytes)[] calls) returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (_Multicall2 *Multicall2TransactorSession) TryBlockAndAggregate(requireSuccess bool, calls []Multicall2Call) (*types.Transaction, error) {
	return _Multicall2.Contract.TryBlockAndAggregate(&_Multicall2.TransactOpts, requireSuccess, calls)
}

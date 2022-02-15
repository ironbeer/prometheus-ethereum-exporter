package contract

import "github.com/umbracle/go-web3"

var (
	LAYER2_PREDEPLOY_CONTRACT_ADDRESSES = map[string]web3.Address{
		"OVM_L2ToL1MessagePasser": web3.HexToAddress("0x4200000000000000000000000000000000000000"),
		"OVM_DeployerWhitelist":   web3.HexToAddress("0x4200000000000000000000000000000000000002"),
		"L2CrossDomainMessenger":  web3.HexToAddress("0x4200000000000000000000000000000000000007"),
		"OVM_GasPriceOracle":      web3.HexToAddress("0x420000000000000000000000000000000000000F"),
		"L2StandardBridge":        web3.HexToAddress("0x4200000000000000000000000000000000000010"),
		"OVM_SequencerFeeVault":   web3.HexToAddress("0x4200000000000000000000000000000000000011"),
		"L2StandardTokenFactory":  web3.HexToAddress("0x4200000000000000000000000000000000000012"),
		"OVM_L1BlockNumber":       web3.HexToAddress("0x4200000000000000000000000000000000000013"),
		"OVM_ETH":                 web3.HexToAddress("0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000"),
		"WETH9":                   web3.HexToAddress("0x4200000000000000000000000000000000000006"),
	}
)

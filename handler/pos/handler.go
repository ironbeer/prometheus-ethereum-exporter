package pos

import (
	"bytes"
	_ "embed"
	"net/http"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"ironbeer/prometheus-ethereum-exporter/contracts/multicall2"
	util "ironbeer/prometheus-ethereum-exporter/handler/internal"
)

var (
	//go:embed abi/stakingmanager.json
	abiJson []byte

	abi *ethabi.ABI
)

func init() {
	if x, err := ethabi.JSON(bytes.NewReader(abiJson)); err != nil {
		panic(err)
	} else {
		abi = &x
	}
}

func Validator(w http.ResponseWriter, r *http.Request) error {
	rpc := r.URL.Query().Get("rpc")
	stakeManagerAddr := common.HexToAddress(r.URL.Query().Get("stakingmanager"))
	validatorAddr := common.HexToAddress(r.URL.Query().Get("validator"))
	mcallAddr := common.HexToAddress(r.URL.Query().Get("multicall"))

	contract := &util.Contract{
		Name: "StakingManager",
		ABI:  abi,
		Methods: []*util.Method{
			{
				Name:   "getValidatorInfo",
				Inputs: []interface{}{validatorAddr, common.Big0},
				Metrics: []*util.Metric{
					{
						Name:     "eth_pos_validator_candidate",
						Output:   "candidate",
						OutputFn: util.BoolOutput,
					},
				},
			},
			{
				Name:   "getBlockAndSlashes",
				Inputs: []interface{}{validatorAddr, common.Big0},
				Metrics: []*util.Metric{
					{
						Name:     "eth_pos_validator_blocks",
						Output:   "blocks",
						OutputFn: util.BigOutput,
					},
					{
						Name:     "eth_pos_validator_slashes",
						Output:   "slashes",
						OutputFn: util.BigOutput,
					},
				},
			},
		},
	}

	calls, err := contract.Calls(stakeManagerAddr)
	if err != nil {
		return err
	}

	client, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}
	defer client.Close()

	mcall, err := multicall2.NewMulticall2(mcallAddr, client)
	if err != nil {
		return err
	}
	result, err := mcall.Aggregate(&bind.CallOpts{Context: r.Context()}, calls)
	if err != nil {
		return err
	}

	collectors, err := contract.Collectors(result.ReturnData)
	if err != nil {
		return err
	}

	registry := prometheus.NewRegistry()
	for _, col := range collectors {
		registry.MustRegister(col)
		defer registry.Unregister(col)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	return nil
}

package basic

import (
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetBalance(w http.ResponseWriter, r *http.Request) error {
	rpc := r.URL.Query().Get("rpc")
	address := r.URL.Query().Get("address")
	if rpc == "" {
		return errors.New("missing parameter: rpc")
	}
	if address == "" {
		return errors.New("missing parameter: address")
	}

	client, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}
	defer client.Close()

	balance, err := client.BalanceAt(r.Context(), common.HexToAddress(address), nil)
	if err != nil {
		return err
	}

	guage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_balance",
		Help: "Displays ethereum token balance(unit: wei)",
	})
	registry := prometheus.NewRegistry()
	registry.MustRegister(guage)
	defer registry.Unregister(guage)

	guage.Set(bigFloat(balance))

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	return nil
}

func GetBlock(w http.ResponseWriter, r *http.Request) error {
	getBlock := func(rpc string, number *big.Int) (*types.Block, error) {
		client, err := ethclient.Dial(rpc)
		if err != nil {
			return nil, err
		}
		defer client.Close()

		return client.BlockByNumber(r.Context(), number)
	}

	rpc := r.URL.Query().Get("rpc")
	origin := r.URL.Query().Get("origin")
	if rpc == "" {
		return errors.New("missing parameter: rpc")
	}

	numberLabel := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("number")))
	if numberLabel == "" {
		numberLabel = "latest"
	}
	labels := prometheus.Labels{"number": numberLabel}

	number := new(big.Int)
	switch numberLabel {
	case "pending":
		number.SetInt64(-1)
	case "latest":
		number.SetInt64(-2)
	case "finalized":
		number.SetInt64(-3)
	case "safe":
		number.SetInt64(-4)
	}

	blk, err := getBlock(rpc, number)
	if err != nil {
		return err
	}

	blockNumber := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "eth_block_number",
		Help:        "Displays ethereum latest block number",
		ConstLabels: labels,
	})
	baseFeePerGas := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "eth_block_baseFeePerGas",
		Help:        "Displays ethereum latest block baseFeePerGas",
		ConstLabels: labels,
	})
	timestamp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "eth_block_timestamp",
		Help:        "Displays ethereum latest block timestamp",
		ConstLabels: labels,
	})
	gasLimit := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "eth_block_gasLimit",
		Help:        "Displays ethereum latest block gas limit",
		ConstLabels: labels,
	})
	gasUsed := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "eth_block_gasUsed",
		Help:        "Displays ethereum latest block gas used",
		ConstLabels: labels,
	})
	transactions := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "eth_block_transactions",
		Help:        "Displays ethereum latest block transaction count",
		ConstLabels: labels,
	})

	collectors := []prometheus.Collector{
		blockNumber,
		baseFeePerGas,
		timestamp,
		gasLimit,
		gasUsed,
		transactions,
	}

	blockNumber.Set(bigFloat(blk.Number()))
	if bf := blk.BaseFee(); bf != nil {
		baseFeePerGas.Set(bigFloat(bf))
	}
	timestamp.Set(float64(blk.Time()))
	gasLimit.Set((float64(blk.GasLimit())))
	gasUsed.Set((float64(blk.GasUsed())))
	transactions.Set((float64(blk.Transactions().Len())))

	if origin != "" {
		originBlk, err := getBlock(origin, blk.Number())
		if err != nil {
			return err
		}

		matchHash := prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "eth_block_match_origin_hash",
			Help:        "Checking if the block hash matches the origin",
			ConstLabels: labels,
		})
		matchState := prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "eth_block_match_origin_state",
			Help:        "Checking if the block state matches the origin",
			ConstLabels: labels,
		})
		collectors = append(collectors, matchHash, matchState)

		if blk.Hash() == originBlk.Hash() {
			matchHash.Set(1)
		}
		if blk.Root() == originBlk.Root() {
			matchState.Set(1)
		}
	}

	registry := prometheus.NewRegistry()
	for _, c := range collectors {
		c := c
		registry.MustRegister(c)
		defer registry.Unregister(c)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	return nil
}

func bigFloat(number *big.Int) float64 {
	fnumber, _ := new(big.Float).SetInt(number).Float64()
	return fnumber
}

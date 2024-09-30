package basic

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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
	rpc := r.URL.Query().Get("rpc")
	if rpc == "" {
		return errors.New("missing parameter: rpc")
	}

	number, numLabel := bigBlockNumber(r.URL.Query().Get("number"))
	labels := make(prometheus.Labels)
	if numLabel != nil {
		labels["number"] = *numLabel
	}

	blk, err := getBlock(r.Context(), rpc, number)
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

func BlockSyncOrigin(w http.ResponseWriter, r *http.Request) error {
	rpc := r.URL.Query().Get("rpc")
	originRPC := r.URL.Query().Get("origin")
	if rpc == "" {
		return errors.New("missing parameter: rpc")
	}
	if originRPC == "" {
		return errors.New("missing parameter: origin")
	}

	number, numLabel := bigBlockNumber(r.URL.Query().Get("number"))
	labels := make(prometheus.Labels)
	if numLabel != nil {
		labels["number"] = *numLabel
	}

	blk, err := getBlock(r.Context(), rpc, number)
	if err != nil {
		return err
	}
	originBlk, err := getBlock(r.Context(), originRPC, blk.Number())
	if err != nil {
		return err
	}

	syncFull := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "eth_block_sync_origin",
		Help:        "Checking if the block hash and state are synchronized with the origin",
		ConstLabels: labels,
	})
	syncHash := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "eth_block_sync_origin_hash",
		Help:        "Checking if the block hash are synchronized with the origin",
		ConstLabels: labels,
	})
	syncRoot := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "eth_block_sync_origin_state",
		Help:        "Checking if the block state are synchronized with the origin",
		ConstLabels: labels,
	})

	collectors := []prometheus.Collector{syncFull, syncHash, syncRoot}

	hashOk := blk.Hash() == originBlk.Hash()
	stateOk := blk.Root() == originBlk.Root()
	if hashOk && stateOk {
		syncFull.Set(1)
	}
	if hashOk {
		syncHash.Set(1)
	}
	if stateOk {
		syncRoot.Set(1)
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

func getBlock(ctx context.Context, rpc string, number *big.Int) (*types.Block, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.BlockByNumber(ctx, number)
}

func bigFloat(number *big.Int) float64 {
	fnumber, _ := new(big.Float).SetInt(number).Float64()
	return fnumber
}

func bigBlockNumber(s string) (num *big.Int, label *string) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "0x") {
		return common.HexToAddress(s).Big(), nil
	}
	if t, err := strconv.ParseInt(s, 10, 64); err == nil {
		return big.NewInt(t), nil
	}

	num = new(big.Int)
	switch s {
	case "earliest":
		num.SetInt64(int64(rpc.EarliestBlockNumber))
	case "pending":
		num.SetInt64(int64(rpc.PendingBlockNumber))
	case "finalized":
		num.SetInt64(int64(rpc.FinalizedBlockNumber))
	case "safe":
		num.SetInt64(int64(rpc.SafeBlockNumber))
	default:
		s = "latest"
		num.SetInt64(int64(rpc.LatestBlockNumber))
	}
	label = &s
	return num, label
}

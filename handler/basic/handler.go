package basic

import (
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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
	rpc := r.URL.Query().Get("rpc")
	if rpc == "" {
		return errors.New("missing parameter: rpc")
	}

	numberLabel := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("number")))
	if numberLabel == "" {
		numberLabel = "latest"
	}

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

	client, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}
	defer client.Close()

	res, err := client.BlockByNumber(r.Context(), number)
	if err != nil {
		return err
	}

	labels := prometheus.Labels{"number": numberLabel}

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

	registry := prometheus.NewRegistry()
	for _, c := range collectors {
		c := c
		registry.MustRegister(c)
		defer registry.Unregister(c)
	}

	blockNumber.Set(bigFloat(res.Number()))
	if bf := res.BaseFee(); bf != nil {
		baseFeePerGas.Set(bigFloat(bf))
	}
	timestamp.Set(float64(res.Time()))
	gasLimit.Set((float64(res.GasLimit())))
	gasUsed.Set((float64(res.GasUsed())))
	transactions.Set((float64(res.Transactions().Len())))

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	return nil
}

func bigFloat(number *big.Int) float64 {
	fnumber, _ := new(big.Float).SetInt(number).Float64()
	return fnumber
}

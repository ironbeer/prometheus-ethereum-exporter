package basic

import (
	"errors"
	"math/big"
	"net/http"

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

	client, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}
	defer client.Close()

	block, err := client.BlockByNumber(r.Context(), nil)
	if err != nil {
		return err
	}

	number := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_block_number",
		Help: "Displays ethereum latest block number",
	})
	baseFeePerGas := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_block_baseFeePerGas",
		Help: "Displays ethereum latest block baseFeePerGas",
	})
	timestamp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_block_timestamp",
		Help: "Displays ethereum latest block timestamp",
	})
	gasLimit := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_block_gasLimit",
		Help: "Displays ethereum latest block gas limit",
	})
	gasUsed := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_block_gasUsed",
		Help: "Displays ethereum latest block gas used",
	})
	transactions := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_block_transactions",
		Help: "Displays ethereum latest block transaction count",
	})

	collectors := []prometheus.Collector{
		number,
		baseFeePerGas,
		timestamp,
		gasLimit,
		gasUsed,
		transactions,
	}

	registry := prometheus.NewRegistry()
	for _, c := range collectors {
		registry.MustRegister(c)
		defer registry.Unregister(c)
	}

	number.Set(bigFloat(block.Number()))
	if bf := block.BaseFee(); bf != nil {
		baseFeePerGas.Set(bigFloat(bf))
	}
	timestamp.Set(float64(block.Time()))
	gasLimit.Set((float64(block.GasLimit())))
	gasUsed.Set((float64(block.GasUsed())))
	transactions.Set((float64(block.Transactions().Len())))

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	return nil
}

func bigFloat(number *big.Int) float64 {
	fnumber, _ := new(big.Float).SetInt(number).Float64()
	return fnumber
}

package main

import (
	"ironbeer/prometheus-ethereum-exporter/optimism"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	web3 "github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	webConfig     = webflag.AddFlags(kingpin.CommandLine)
	listenAddress = kingpin.Flag("web.listen", "The address to listen on for HTTP requests.").Default(":49000").String()
	abiPath       = kingpin.Flag("optimism.abi", "Optimism L1 ABI directory path").Default("").String()

	methods = map[string]func(w http.ResponseWriter, r *http.Request){
		"eth.getBalance":  getBalance,
		"eth.getBlock":    getBlock,
		"optimism.status": optimismStatus,
	}
)

func bigToFloat(number *big.Int) float64 {
	fnumber, _ := new(big.Float).SetInt(number).Float64()
	return fnumber
}

func parseUint64orHex(str string) (uint64, error) {
	base := 10
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
		base = 16
	}
	return strconv.ParseUint(str, base, 64)
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	rpc := r.URL.Query().Get("rpc")
	address := r.URL.Query().Get("address")
	if rpc == "" || address == "" {
		return
	}

	client, err := jsonrpc.NewClient(rpc)
	if err != nil {
		return
	}

	number, err := client.Eth().GetBalance(web3.HexToAddress(address), web3.Latest)
	if err != nil {
		return
	}

	guage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_balance",
		Help: "Displays ethereum token balance(unit: wei)",
	})
	registry := prometheus.NewRegistry()
	registry.MustRegister(guage)

	guage.Set(bigToFloat(number))

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func getBlock(w http.ResponseWriter, r *http.Request) {
	rpc := r.URL.Query().Get("rpc")
	if rpc == "" {
		return
	}

	client, err := jsonrpc.NewClient(rpc)
	if err != nil {
		return
	}

	block, err := client.Eth().GetBlockByNumber(web3.Latest, false)
	if err != nil {
		return
	}

	var out string
	if err := client.Call("eth_getBlockTransactionCountByNumber", &out, "pending"); err != nil {
		return
	}
	pendingTx, err := parseUint64orHex(out)
	if err != nil {
		return
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
	pendingTransactions := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_block_pendingTransactions",
		Help: "Displays ethereum latest pending transaction count",
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(number)
	registry.MustRegister(baseFeePerGas)
	registry.MustRegister(timestamp)
	registry.MustRegister(gasLimit)
	registry.MustRegister(gasUsed)
	registry.MustRegister(transactions)
	registry.MustRegister(pendingTransactions)

	number.Set(float64(block.Number))
	baseFeePerGas.Set(float64(block.BaseFeePerGas))
	timestamp.Set(float64(block.Timestamp))
	gasLimit.Set(float64(block.GasLimit))
	gasUsed.Set(float64(block.GasUsed))
	transactions.Set(float64(len(block.TransactionsHashes)))
	pendingTransactions.Set(float64(pendingTx))

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func optimismStatus(w http.ResponseWriter, r *http.Request) {
	rpc := r.URL.Query().Get("rpc")
	ctc := r.URL.Query().Get("ctc")
	if rpc == "" || ctc == "" {
		return
	}

	client, err := jsonrpc.NewClient(rpc)
	if err != nil {
		return
	}

	oclient := optimism.NewOptimismClient(client, *abiPath)
	contract, err := oclient.GetContract("CanonicalTransactionChain", web3.HexToAddress(ctc))
	if err != nil {
		return
	}

	totalElements := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_optimism_ctc_total_elements",
		Help: " Retrieves the total number of elements submitted",
	})
	totalBatches := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_optimism_ctc_total_batches",
		Help: "Retrieves the total number of batches submitted",
	})
	nextQueueIndex := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_optimism_ctc_next_queue_index",
		Help: "Returns the index of the next element to be enqueued",
	})
	lastTimestamp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_optimism_ctc_last_timestamp",
		Help: "Returns the timestamp of the last transaction",
	})
	lastBlockNumber := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_optimism_ctc_last_block_number",
		Help: "Returns the blocknumber of the last transaction",
	})
	numPendingQueueElements := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_optimism_ctc_num_pending_queue_elements",
		Help: "Get the number of queue elements which have not yet been included",
	})
	queueLength := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "eth_optimism_ctc_queue_length",
		Help: "Retrieves the length of the queue, including both pending and canonical transactions",
	})

	metrics := map[string]struct {
		gauge prometheus.Gauge
		key   string
	}{
		"getTotalElements":           {totalElements, "_totalElements"},
		"getTotalBatches":            {totalBatches, "_totalBatches"},
		"getNextQueueIndex":          {nextQueueIndex, "0"},
		"getLastTimestamp":           {lastTimestamp, "0"},
		"getLastBlockNumber":         {lastBlockNumber, "0"},
		"getNumPendingQueueElements": {numPendingQueueElements, "0"},
		"getQueueLength":             {queueLength, "0"},
	}

	registry := prometheus.NewRegistry()

	var wg sync.WaitGroup
	wg.Add(len(metrics))

	for method, m := range metrics {
		go func(gauge prometheus.Gauge, method, key string) {
			defer wg.Done()

			registry.MustRegister(gauge)

			result, err := contract.Call(method, web3.Latest)
			if err != nil {
				return
			}

			val, ok := oclient.DecodeBigInt(result, key)
			if !ok {
				return
			}

			gauge.Set(val)
		}(m.gauge, method, m.key)
	}

	wg.Wait()

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		method, ok := methods[r.URL.Query().Get("method")]
		if !ok {
			return
		}
		method(w, r)
	})

	srv := &http.Server{Addr: *listenAddress}
	srvc := make(chan struct{})
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := web.ListenAndServe(srv, *webConfig, logger); err != http.ErrServerClosed {
			close(srvc)
		}
	}()

	for {
		select {
		case <-term:
			level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
			os.Exit(0)
		case <-srvc:
			os.Exit(1)
		}
	}
}

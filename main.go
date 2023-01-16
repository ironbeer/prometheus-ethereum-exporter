package main

import (
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	web3 "github.com/umbracle/go-web3"
	web3contract "github.com/umbracle/go-web3/contract"
	"github.com/umbracle/go-web3/jsonrpc"
	"gopkg.in/alecthomas/kingpin.v2"

	"ironbeer/prometheus-ethereum-exporter/contract"
)

var (
	webConfig     = webflag.AddFlags(kingpin.CommandLine)
	listenAddress = kingpin.
			Flag("web.listen", "The address to listen on for HTTP requests.").
			Default(":49000").
			String()
	abiPath = kingpin.Flag("abi", "Contract ABI directory path").Default("").String()

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
	defer client.Close()

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
	defer registry.Unregister(guage)

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
	defer client.Close()

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

	collectors := []prometheus.Collector{
		number,
		baseFeePerGas,
		timestamp,
		gasLimit,
		gasUsed,
		transactions,
		pendingTransactions,
	}

	registry := prometheus.NewRegistry()
	for _, c := range collectors {
		registry.MustRegister(c)
		defer registry.Unregister(c)
	}

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
	rpcUrl := r.URL.Query().Get("rpc")
	ctcAddress := r.URL.Query().Get("ctc")
	sccAddress := r.URL.Query().Get("scc")
	if rpcUrl == "" || ctcAddress == "" || sccAddress == "" {
		return
	}

	client, err := jsonrpc.NewClient(rpcUrl)
	if err != nil {
		return
	}
	defer client.Close()

	ctc, err := contract.GetContract(
		client, web3.HexToAddress(ctcAddress),
		filepath.Join(*abiPath, "CanonicalTransactionChain.json"))
	if err != nil {
		return
	}

	scc, err := contract.GetContract(
		client, web3.HexToAddress(sccAddress),
		filepath.Join(*abiPath, "OasysStateCommitmentChain.json"))
	if err != nil {
		return
	}

	metrics := []struct {
		contract    *web3contract.Contract
		method, key string
		collector   prometheus.Gauge
	}{
		// CanonicalTransactionChain
		{
			ctc, "getTotalElements", "_totalElements",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_ctc_total_elements",
				Help: " Retrieves the total number of elements submitted",
			}),
		},
		{
			ctc, "getTotalBatches", "_totalBatches",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_ctc_total_batches",
				Help: "Retrieves the total number of batches submitted",
			}),
		},
		{
			ctc, "getNextQueueIndex", "0",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_ctc_next_queue_index",
				Help: "Returns the index of the next element to be enqueued",
			}),
		},
		{
			ctc, "getLastTimestamp", "0",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_ctc_last_timestamp",
				Help: "Returns the timestamp of the last transaction",
			}),
		},
		{
			ctc, "getLastBlockNumber", "0",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_ctc_last_block_number",
				Help: "Returns the blocknumber of the last transaction",
			}),
		},
		{
			ctc, "getNumPendingQueueElements", "0",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_ctc_num_pending_queue_elements",
				Help: "Get the number of queue elements which have not yet been included",
			}),
		},
		{
			ctc, "getQueueLength", "0",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_ctc_queue_length",
				Help: "Retrieves the length of the queue, including both pending and canonical transactions",
			}),
		},
		// StateCommitmentChain
		{
			scc, "getTotalElements", "_totalElements",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_scc_total_elements",
				Help: " Retrieves the total number of elements submitted",
			}),
		},
		{
			scc, "getTotalBatches", "_totalBatches",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_scc_total_batches",
				Help: "Retrieves the total number of batches submitted",
			}),
		},
		{
			scc, "getLastSequencerTimestamp", "_lastSequencerTimestamp",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_scc_last_sequencer_timestamp",
				Help: "Retrieves the timestamp of the last batch submitted by the sequencer",
			}),
		},
		{
			scc, "nextIndex", "0",
			prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "eth_optimism_scc_next_index",
				Help: "Retrieves the batch index to verify next",
			}),
		},
	}

	registry := prometheus.NewRegistry()

	for _, m := range metrics {
		registry.MustRegister(m.collector)
		defer registry.Unregister(m.collector)

		result, err := m.contract.Call(m.method, web3.Latest)
		if err != nil {
			fmt.Printf("method: %s, key: %s, err: %s\n", m.method, m.key, err.Error())
			return
		}

		val, ok := contract.DecodeBigInt(result, m.key)
		if !ok {
			return
		}

		m.collector.Set(val)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if method, ok := methods[r.URL.Query().Get("method")]; ok {
			method(w, r)
		}
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

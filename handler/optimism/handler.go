package optimism

import (
	"bytes"
	_ "embed"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"ironbeer/prometheus-ethereum-exporter/contracts/multicall2"
	util "ironbeer/prometheus-ethereum-exporter/handler/internal"
)

var (
	//go:embed abi/ctc.json
	ctcABI []byte
	//go:embed abi/scc.json
	sccABI []byte
	//go:embed abi/l2oo.json
	l2ooABI []byte

	contracts = []struct {
		query    string
		abi      *[]byte
		contract *util.Contract
	}{
		{
			query: "ctc",
			abi:   &ctcABI,
			contract: &util.Contract{
				Name: "CTC",
				Methods: []*util.Method{
					{
						Name: "getTotalElements",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_ctc_total_elements",
								Help:   "Retrieves the total number of elements submitted",
								Output: util.BigOutput("_totalElements"),
							},
						},
					},
					{
						Name: "getTotalBatches",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_ctc_total_batches",
								Help:   "Retrieves the total number of batches submitted",
								Output: util.BigOutput("_totalBatches"),
							},
						},
					},
					{
						Name: "getNextQueueIndex",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_ctc_next_queue_index",
								Help:   "Returns the index of the next element to be enqueued",
								Output: util.BigOutput(""),
							},
						},
					},
					{
						Name: "getLastTimestamp",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_ctc_last_timestamp",
								Help:   "Returns the timestamp of the last transaction",
								Output: util.BigOutput(""),
							},
						},
					},
					{
						Name: "getLastBlockNumber",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_ctc_last_block_number",
								Help:   "Returns the blocknumber of the last transaction",
								Output: util.BigOutput(""),
							},
						},
					},
					{
						Name: "getNumPendingQueueElements",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_ctc_num_pending_queue_elements",
								Help:   "Get the number of queue elements which have not yet been included",
								Output: util.BigOutput(""),
							},
						},
					},
					{
						Name: "getQueueLength",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_ctc_queue_length",
								Help:   "Retrieves the length of the queue, including both pending and canonical transactions",
								Output: util.BigOutput(""),
							},
						},
					},
				},
			},
		},
		{
			query: "scc",
			abi:   &sccABI,
			contract: &util.Contract{
				Name: "SCC",
				Methods: []*util.Method{
					{
						Name: "getTotalElements",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_scc_total_elements",
								Help:   "Retrieves the total number of elements submitted",
								Output: util.BigOutput("_totalElements"),
							},
						},
					},
					{
						Name: "getTotalBatches",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_scc_total_batches",
								Help:   "Retrieves the total number of batches submitted",
								Output: util.BigOutput("_totalBatches"),
							},
						},
					},
					{
						Name: "getLastSequencerTimestamp",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_scc_last_sequencer_timestamp",
								Help:   "Retrieves the timestamp of the last batch submitted by the sequencer",
								Output: util.BigOutput("_lastSequencerTimestamp"),
							},
						},
					},
					{
						Name: "nextIndex",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_scc_next_index",
								Help:   "Retrieves the batch index to verify next",
								Output: util.BigOutput(""),
							},
						},
					},
				},
			},
		},
		{
			query: "l2oo",
			abi:   &l2ooABI,
			contract: &util.Contract{
				Name: "L2OO",
				Methods: []*util.Method{
					{
						Name: "latestOutputIndex",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_l2oo_latest_output_index",
								Help:   "Returns the number of outputs that have been proposed",
								Output: util.BigOutput(""),
							},
						},
					},
					{
						Name: "nextVerifyIndex",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_l2oo_next_verify_index",
								Help:   "Next L2Output index to verify",
								Output: util.BigOutput(""),
							},
						},
					},
					{
						Name: "latestBlockNumber",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_l2oo_latest_block_number",
								Help:   "Returns the block number of the latest submitted L2 output proposal",
								Output: util.BigOutput(""),
							},
						},
					},
				},
			},
		},
		{
			query: "opl2oo",
			abi:   &l2ooABI,
			contract: &util.Contract{
				Name: "OPL2OO",
				Methods: []*util.Method{
					{
						Name: "latestOutputIndex",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_l2oo_latest_output_index",
								Help:   "Returns the number of outputs that have been proposed",
								Output: util.BigOutput(""),
							},
						},
					},
					{
						Name: "latestBlockNumber",
						Metrics: []*util.Metric{
							{
								Name:   "eth_optimism_l2oo_latest_block_number",
								Help:   "Returns the block number of the latest submitted L2 output proposal",
								Output: util.BigOutput(""),
							},
						},
					},
				},
			},
		},
	}
)

func init() {
	for _, c := range contracts {
		abi, err := abi.JSON(bytes.NewReader(*c.abi))
		if err != nil {
			panic(err)
		}
		c.contract.ABI = &abi
	}
}

func Status(w http.ResponseWriter, r *http.Request) error {
	rpc := r.URL.Query().Get("rpc")
	mcallAddr := common.HexToAddress(r.URL.Query().Get("multicall"))

	var (
		called []*util.Contract
		mcalls []multicall2.Multicall2Call
	)
	for _, c := range contracts {
		addr := common.HexToAddress(r.URL.Query().Get(c.query))
		if addr == (common.Address{}) {
			continue
		}

		if x, err := c.contract.Calls(addr); err != nil {
			return err
		} else {
			called = append(called, c.contract)
			mcalls = append(mcalls, x...)
		}
	}

	registry := prometheus.NewRegistry()
	response := func() error {
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
		return nil
	}
	if len(mcalls) == 0 {
		return response()
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
	result, err := mcall.Aggregate(&bind.CallOpts{Context: r.Context()}, mcalls)
	if err != nil {
		return err
	}

	var offset, limit int
	for _, c := range called {
		offset = limit
		limit = offset + len(c.Methods)

		collectors, err := c.Collectors(result.ReturnData[offset:limit])
		if err != nil {
			return err
		}

		for _, col := range collectors {
			registry.MustRegister(col)
			defer registry.Unregister(col)
		}
	}

	return response()
}

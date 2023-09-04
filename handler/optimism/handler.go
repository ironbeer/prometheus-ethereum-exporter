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
	ctcABIJson []byte

	//go:embed abi/scc.json
	sccABIJson []byte

	ctcABI, sccABI *abi.ABI
)

func init() {
	if x, err := abi.JSON(bytes.NewReader(ctcABIJson)); err != nil {
		panic(err)
	} else {
		ctcABI = &x
	}

	if x, err := abi.JSON(bytes.NewReader(sccABIJson)); err != nil {
		panic(err)
	} else {
		sccABI = &x
	}
}

func Status(w http.ResponseWriter, r *http.Request) error {
	rpc := r.URL.Query().Get("rpc")
	ctcAddr := common.HexToAddress(r.URL.Query().Get("ctc"))
	sccAddr := common.HexToAddress(r.URL.Query().Get("scc"))
	mcallAddr := common.HexToAddress(r.URL.Query().Get("multicall"))

	ctc := &util.Contract{
		Name: "CTC",
		ABI:  ctcABI,
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
	}

	scc := &util.Contract{
		Name: "SCC",
		ABI:  sccABI,
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
	}

	contracts := []struct {
		contract *util.Contract
		addr     common.Address
	}{
		{ctc, ctcAddr},
		{scc, sccAddr},
	}

	var (
		calls  []multicall2.Multicall2Call
		limits []int
	)
	for _, c := range contracts {
		if x, err := c.contract.Calls(c.addr); err != nil {
			return err
		} else {
			calls = append(calls, x...)
			limits = append(limits, len(x))
		}
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

	registry := prometheus.NewRegistry()

	var offset, limit int
	for i, c := range contracts {
		offset = limit
		limit = limits[i]
		datas := result.ReturnData[offset : offset+limit]

		collectors, err := c.contract.Collectors(datas)
		if err != nil {
			return err
		}

		for _, col := range collectors {
			registry.MustRegister(col)
			defer registry.Unregister(col)
		}
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	return nil
}

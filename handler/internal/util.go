package util

import (
	"errors"
	"fmt"
	"ironbeer/prometheus-ethereum-exporter/contracts/multicall2"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	errTypeMismatch = errors.New("type mismatch")
)

type Methodx = string

type Contract struct {
	Name    string
	ABI     *abi.ABI
	Methods []*Method
}

type Method struct {
	Name    string
	Inputs  []interface{}
	Metrics []*Metric
}

type Metric struct {
	Name, Help string
	Output     string
	OutputFn   OutputFn
}

type OutputFn func(src interface{}) (float64, error)

func (c *Contract) Calls(addr common.Address) (calls []multicall2.Multicall2Call, err error) {
	for _, method := range c.Methods {
		packed, err := c.ABI.Methods[method.Name].Inputs.Pack(method.Inputs...)
		if err != nil {
			return nil, err
		}

		calls = append(calls, multicall2.Multicall2Call{
			Target:   addr,
			CallData: append(c.ABI.Methods[method.Name].ID, packed...),
		})
	}
	return calls, nil
}

func (c *Contract) Collectors(datas [][]byte) ([]prometheus.Collector, error) {
	var collectors []prometheus.Collector

	for i, method := range c.Methods {
		values := map[string]interface{}{}
		if err := c.ABI.Methods[method.Name].Outputs.UnpackIntoMap(values, datas[i]); err != nil {
			return nil, err
		}

		for _, metric := range method.Metrics {
			val, err := metric.OutputFn(values[metric.Output])
			if err != nil {
				return nil, fmt.Errorf(
					"not found: contract=%s method=%s output=%s values=%v",
					c.Name, method.Name, metric.Output, values)
			}

			collector := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: metric.Name,
				Help: metric.Help,
			})
			collector.Set(val)

			collectors = append(collectors, collector)
		}
	}

	return collectors, nil
}

func BigOutput(src interface{}) (float64, error) {
	t, ok := src.(*big.Int)
	if !ok {
		return 0, errTypeMismatch
	}
	return float64(t.Uint64()), nil
}

func BoolOutput(src interface{}) (float64, error) {
	t, ok := src.(bool)
	if !ok {
		return 0, errTypeMismatch
	}
	if t {
		return 1, nil
	}
	return 0, nil
}

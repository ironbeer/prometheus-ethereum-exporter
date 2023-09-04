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
	errValueNotFound = errors.New("value not found")
	errTypeMismatch  = errors.New("type mismatch")
)

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
	Output     OutputFn
}

type OutputFn func(values map[string]interface{}) (float64, error)

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
			val, err := metric.Output(values)
			if err != nil {
				return nil, fmt.Errorf(
					"error: reason=%s contract=%s method=%s metric=%s values=%v",
					err, c.Name, method.Name, metric.Name, values)
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

func BigOutput(output string) OutputFn {
	return func(values map[string]interface{}) (float64, error) {
		if val, ok := values[output]; !ok {
			return 0, errValueNotFound
		} else if t, ok := val.(*big.Int); !ok {
			return 0, errTypeMismatch
		} else {
			return float64(t.Uint64()), nil
		}
	}
}

func BoolOutput(output string) OutputFn {
	return func(values map[string]interface{}) (float64, error) {
		if val, ok := values[output]; !ok {
			return 0, errValueNotFound
		} else if t, ok := val.(bool); !ok {
			return 0, errTypeMismatch
		} else if t {
			return 1, nil
		}
		return 0, nil
	}
}

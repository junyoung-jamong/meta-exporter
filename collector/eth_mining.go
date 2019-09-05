package collector

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthMining struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

func NewEthMining(rpc *rpc.Client) *EthMining {
	return &EthMining{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"meta_mining",
			"Returns 1 if client is actively mining new blocks.",
			nil,
			nil,
		),
	}
}

func (collector *EthMining) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthMining) Collect(ch chan<- prometheus.Metric) {
	var result bool
	if err := collector.rpc.Call(&result, "eth_mining"); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	if result {
		value := float64(1)
		ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
	} else {
		value := float64(0)
		ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
	}
}

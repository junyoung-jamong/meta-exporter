package collector

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type NetListening struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

func NewNetListening(rpc *rpc.Client) *NetListening {
	return &NetListening{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"net_listening",
			"returns 1 if client is actively listening for network connections",
			nil,
			nil,
		),
	}
}

func (collector *NetListening) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *NetListening) Collect(ch chan<- prometheus.Metric) {
	var result bool
	if err := collector.rpc.Call(&result, "net_listening"); err != nil {
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

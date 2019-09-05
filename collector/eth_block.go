package collector

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/onrik/ethrpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthBlock struct {
	rpc             *rpc.Client
	url             string
	blockNumberDesc *prometheus.Desc
	timeStampDesc   *prometheus.Desc
	gasUsedDesc     *prometheus.Desc
	gasLimitDesc    *prometheus.Desc
	blockSizeDesc   *prometheus.Desc
}

func NewEthBlock(rpc *rpc.Client, url string) *EthBlock {
	return &EthBlock{
		url: url,
		rpc: rpc,
		blockNumberDesc: prometheus.NewDesc(
			"block_number",
			"the number of most recent block.",
			nil,
			nil,
		),
		timeStampDesc: prometheus.NewDesc(
			"block_timestamp",
			"the unix timestamp for when the block was collated.",
			nil,
			nil,
		),
		gasUsedDesc: prometheus.NewDesc(
			"block_gas_used",
			"the total used gas by all transactions in this block.",
			nil,
			nil,
		),
		gasLimitDesc: prometheus.NewDesc(
			"block_limit",
			"the maximum gas allowed in this block.",
			nil,
			nil,
		),
		blockSizeDesc: prometheus.NewDesc(
			"block_size",
			"integer the size of this block in bytes.",
			nil,
			nil,
		),
	}
}

func (collector *EthBlock) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.blockNumberDesc
	ch <- collector.timeStampDesc
	ch <- collector.gasUsedDesc
	ch <- collector.gasLimitDesc
	ch <- collector.blockSizeDesc
}

func (collector *EthBlock) Collect(ch chan<- prometheus.Metric) {
	client := ethrpc.New(collector.url)

	number, err := client.EthBlockNumber()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.blockNumberDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.timeStampDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.gasUsedDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.gasLimitDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.blockSizeDesc, err)
		return
	}
	//block, err := client.getBlock("eth_getBlockByNumber", true, "latest", true)
	block, err := client.EthGetBlockByNumber(number, true)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.blockNumberDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.timeStampDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.gasUsedDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.gasLimitDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.blockSizeDesc, err)
		return
	}

	value := float64(number)
	ch <- prometheus.MustNewConstMetric(collector.blockNumberDesc, prometheus.GaugeValue, value)
	value = float64(block.Timestamp)
	ch <- prometheus.MustNewConstMetric(collector.timeStampDesc, prometheus.GaugeValue, value)
	value = float64(block.GasUsed)
	ch <- prometheus.MustNewConstMetric(collector.gasUsedDesc, prometheus.GaugeValue, value)
	value = float64(block.GasLimit)
	ch <- prometheus.MustNewConstMetric(collector.gasLimitDesc, prometheus.GaugeValue, value)
	value = float64(block.Size)
	ch <- prometheus.MustNewConstMetric(collector.blockSizeDesc, prometheus.GaugeValue, value)
}

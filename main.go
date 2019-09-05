package main

import (
	"flag"
	"fmt"
	"log"
	"meta-exporter/collector"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	flag.Usage = func() {
		const (
			usage = "Usage: metadium_exporter [options]\n\n" +
				"Prometheus exporter for Metadium client metrics\n\n" +
				"Options:\n"
		)

		fmt.Fprint(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()

		os.Exit(2)
	}

	url := flag.String("url", "http://localhost:8588", "Metadium JSON-RPC URL")
	addr := flag.String("addr", ":9101", "listen address")

	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
	}

	rpc, err := rpc.Dial(*url)
	if err != nil {
		log.Fatal(err)
	}

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(
		collector.NewNetPeerCount(rpc),
		collector.NewEthGasPrice(rpc),
		collector.NewEthEarliestBlockTransactions(rpc),
		collector.NewEthLatestBlockTransactions(rpc),
		collector.NewEthPendingBlockTransactions(rpc),
		collector.NewNetListening(rpc),
		collector.NewEthMining(rpc),
		collector.NewEthBlock(rpc, *url),
	)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.New(os.Stderr, log.Prefix(), log.Flags()),
		ErrorHandling: promhttp.ContinueOnError,
	})

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

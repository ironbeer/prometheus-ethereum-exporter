package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"gopkg.in/alecthomas/kingpin.v2"

	"ironbeer/prometheus-ethereum-exporter/handler/basic"
	"ironbeer/prometheus-ethereum-exporter/handler/optimism"
	"ironbeer/prometheus-ethereum-exporter/handler/pos"
)

var (
	webConfig     = webflag.AddFlags(kingpin.CommandLine)
	listenAddress = kingpin.
			Flag("web.listen", "The address to listen on for HTTP requests.").
			Default(":49000").
			String()

	methods = map[string]func(http.ResponseWriter, *http.Request) error{
		"eth.getBalance":      basic.GetBalance,
		"eth.getBlock":        basic.GetBlock,
		"eth.blockSyncOrigin": basic.BlockSyncOrigin,
		"pos.validator":       pos.Validator,
		"optimism.status":     optimism.Status,
	}
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("method")
		if method, ok := methods[name]; !ok {
			http.Error(w, fmt.Sprintf("unknown method: %s", name), http.StatusBadRequest)
		} else if err := method(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

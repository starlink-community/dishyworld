package agent

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/scrape"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/prometheus/prometheus/tsdb"
	"os"

	"time"
)

type MiniProm struct {
	logger     log.Logger
	conf       *config.Config
	tsdbDir    string
	addr       string
	remoteAddr string
}

func NewMiniProm(tsdbDir string, addr string, remoteAddr string, dishId string) (*MiniProm, error) {
	ll := &promlog.AllowedLevel{}
	ll.Set("warn")
	s := fmt.Sprintf(`
scrape_configs:
  - job_name: 'starlink'
    static_configs:
    - targets: ['%s']

    # this allows us to use the dish-id as a target, but still poll localhost
    relabel_configs:
    - source_labels: [__address__]
      target_label: __param_target
    - source_labels: [__param_target]
      target_label: instance
    - target_label: __address__
      replacement: %s

# yes I know I am writing credentials into git. Will clean this up once a few folks have tried it
remote_write:
  - url: https://prometheus-us-central1.grafana.net/api/prom/push
    basic_auth:
      username: 44690
      password: eyJrIjoiMzEyNTAwMzI1NTNlOWU5ZTY2ZDcxZDA5ZjhjYWM1MmMxZTY0MzIzMCIsIm4iOiJjbGllbnQiLCJpZCI6NDY1NDAzfQ==

`, dishId, addr)

	cfg, err := config.Load(s)
	if err != nil {
		return nil, err
	}
	return &MiniProm{
		tsdbDir:    tsdbDir,
		logger:     promlog.New(&promlog.Config{Level: ll}),
		conf:       cfg,
		addr:       addr,
		remoteAddr: remoteAddr,
	}, nil
}
func (m *MiniProm) Start() {
	db, err := tsdb.Open(
		m.tsdbDir,
		log.With(m.logger, "component", "tsdb"),
		prometheus.DefaultRegisterer,
		nil,
	)
	if err != nil {
		fmt.Println("db failed", err)
		os.Exit(1)
	}
	var (
		localStorage           = db
		remoteStorage          = remote.NewStorage(log.With(m.logger, "component", "remote"), prometheus.DefaultRegisterer, localStorage.StartTime, m.tsdbDir, time.Duration(1*time.Minute), nil)
		fanoutStorage          = storage.NewFanout(m.logger, localStorage, remoteStorage)
		scrapeManager          = scrape.NewManager(log.With(m.logger, "component", "scrape manager"), fanoutStorage)
		ctxScrape, _           = context.WithCancel(context.Background())
		discoveryManagerScrape = discovery.NewManager(ctxScrape, log.With(m.logger, "component", "discovery manager scrape"), discovery.Name("scrape"))
	)
	scrapeManager.ApplyConfig(m.conf)
	remoteStorage.ApplyConfig(m.conf)
	c := make(map[string]discovery.Configs)
	for _, v := range m.conf.ScrapeConfigs {
		c[v.JobName] = v.ServiceDiscoveryConfigs
	}
	discoveryManagerScrape.ApplyConfig(c)
	go func() {
		discoveryManagerScrape.Run()
	}()

	go func() {
		scrapeManager.Run(discoveryManagerScrape.SyncCh())
	}()
}

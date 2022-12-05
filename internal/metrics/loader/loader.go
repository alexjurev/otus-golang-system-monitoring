package metricloader

import (
	"errors"
	"fmt"

	"github.com/alexjurev/otus-golang-system-monitoring/internal/metrics"
	"github.com/alexjurev/otus-golang-system-monitoring/internal/metrics/cpu"
	"github.com/alexjurev/otus-golang-system-monitoring/internal/metrics/diskinfo"
	"github.com/alexjurev/otus-golang-system-monitoring/internal/metrics/loadaverage"
	"github.com/alexjurev/otus-golang-system-monitoring/internal/metrics/spaceinfo"
)

type Config struct {
	IgnoreUnavailable bool
	Collect           Metric
}

type Metric struct {
	Cpu         bool
	LoadAverage bool
	DiskInfo    bool
	SpaceInfo   bool
}

var ErrCollectorNotAvailable = errors.New("collector is not available")

func Load(config Config) ([]metric.Collector, error) {
	var collectors []metric.Collector
	var err error
	if config.Collect.Cpu {
		collectors, err = appendCollector(collectors, cpu.Collector{}, config.IgnoreUnavailable)
		if err != nil {
			return nil, err
		}
	}
	if config.Collect.LoadAverage {
		collectors, err = appendCollector(collectors, loadaverage.Collector{}, config.IgnoreUnavailable)
		if err != nil {
			return nil, err
		}
	}
	if config.Collect.DiskInfo {
		collectors, err = appendCollector(collectors, diskinfo.Collector{}, config.IgnoreUnavailable)
		if err != nil {
			return nil, err
		}
	}
	if config.Collect.SpaceInfo {
		collectors, err = appendCollector(collectors, spaceinfo.Collector{}, config.IgnoreUnavailable)
		if err != nil {
			return nil, err
		}
	}

	return collectors, nil
}

func appendCollector(collectors []metric.Collector, collector metric.Collector, ignoreNotAvailable bool) ([]metric.Collector, error) {
	if !collector.Available() {
		if !ignoreNotAvailable {
			return collectors, fmt.Errorf("failed load collector '%s': %w", collector.Name(), ErrCollectorNotAvailable)
		}
		return collectors, nil
	}
	return append(collectors, collector), nil
}

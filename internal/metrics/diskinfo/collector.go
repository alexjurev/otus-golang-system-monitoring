package diskinfo

import (
	metric "github.com/alexjurev/otus-golang-system-monitoring/internal/metrics"
	"time"
)

const (
	collectorName    = "Disk Information"
	groupName        = "DiskInfo"
	tpsMetricName    = "tps"
	rSpeedMetricName = "Blk_read/s"
	wSpeedMetricName = "Blk_wrtn/s"
)

//nolint:deadcode,unused // ignore when collector is not available
type metricData struct {
	TPS    float64
	RSpeed float64
	WSpeed float64
}

/* eslint-enable no-unused-vars */

type Collector struct {
	metric.UnavailableCollector
	// _ metricData // Just for `unused` linter
}

func (c Collector) Name() string {
	return collectorName
}

func (c Collector) GroupName() string {
	return groupName
}

//nolint:deadcode,unused // ignore when collector is not available
func toGroup(t time.Time, m metricData) metric.Group {
	return metric.Group{
		Name: groupName,
		Time: t,
		Metrics: []metric.Metric{
			{
				Time:  t,
				Name:  tpsMetricName,
				Value: m.TPS,
			},
			{
				Time:  t,
				Name:  rSpeedMetricName,
				Value: m.RSpeed,
			},
			{
				Time:  t,
				Name:  wSpeedMetricName,
				Value: m.WSpeed,
			},
		}}
}

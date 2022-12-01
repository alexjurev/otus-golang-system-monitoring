package cpu

import (
	metric "github.com/alexjurev/otus-golang-system-monitoring/internal/metrics"
	"time"
)

const (
	collectorName        = "CPU"
	groupName            = "cpu"
	systemTimeMetricName = "system-time"
	userTimeMetricName   = "user-time"
	idleTimeMetricName   = "idle-time"
)

type metricData struct { //nolint
	SystemTime float64
	UserTime   float64
	IdleTime   float64
}

type Collector struct {
	metric.UnavailableCollector
}

func (c Collector) Name() string {
	return collectorName
}

func (c Collector) GroupName() string {
	return groupName
}

func toGroup(t time.Time, m metricData) metric.Group { //nolint
	return metric.Group{
		Name: groupName,
		Time: t,
		Metrics: []metric.Metric{
			{
				Time:  t,
				Name:  systemTimeMetricName,
				Value: m.SystemTime,
			},
			{
				Time:  t,
				Name:  userTimeMetricName,
				Value: m.UserTime,
			},
			{
				Time:  t,
				Name:  idleTimeMetricName,
				Value: m.IdleTime,
			},
		}}
}

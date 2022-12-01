package spaceinfo

import (
	metric "github.com/alexjurev/otus-golang-system-monitoring/internal/metrics"
	"time"
)

const (
	collectorName              = "Free/Used Space Info"
	groupName                  = "SpaceInfo"
	usedSpaceMetricName        = "1K-blocks"
	percentUsedMetricName      = "Used"
	usedInodeSpaceMetricName   = "Available"
	percentInodeUsedMetricName = "Use%"
)

//nolint:deadcode,unused // ignore when collector is not available
type metricData struct {
	usedSpace     float64
	percentSpace  float64
	usedISpace    float64
	percentIspace float64
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
				Name:  usedSpaceMetricName,
				Value: m.usedSpace,
			},
			{
				Time:  t,
				Name:  percentUsedMetricName,
				Value: m.percentSpace,
			},
			{
				Time:  t,
				Name:  usedInodeSpaceMetricName,
				Value: m.usedISpace,
			},
			{
				Time:  t,
				Name:  percentInodeUsedMetricName,
				Value: m.percentIspace,
			},
		}}
}

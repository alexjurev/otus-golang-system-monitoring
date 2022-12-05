package diskinfo_test

import (
	"github.com/alexjurev/otus-golang-system-monitoring/internal/metrics/diskinfo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetDiskInfoMetric(t *testing.T) {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		var a []int
		var i int
		ch <- struct{}{}
		for {
			i++
			a = append(a, i)
			_ = a
			select {
			case <-ch:
				return
			default:
			}
		}
	}()
	<-ch
	metric, err := diskinfo.Collector{}.GetMetrics()
	ch <- struct{}{}
	require.NoError(t, err)
	require.Greater(t, metric.Metrics[0].Value+metric.Metrics[1].Value+metric.Metrics[2].Value, float64(0))
	<-ch
}

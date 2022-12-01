package diskinfo

import (
	"fmt"
	"github.com/alexjurev/otus-golang-system-monitoring/internal/executor"
	"github.com/alexjurev/otus-golang-system-monitoring/internal/metrics"
	"strconv"
	"strings"
	"time"
	"unicode" //nolint
)

const (
	command = "iostat"
)

func (c Collector) Available() bool {
	return true
}

func (c Collector) GetMetrics() (metric.Group, error) {
	output, err := executor.Exec(command, nil)
	if err != nil {
		return metric.Group{}, err
	}

	var m metricData
	var t time.Time = time.Now()
	if err := parse(output, &m); err != nil {
		return metric.Group{}, err
	}
	return toGroup(t, m), nil
}

func parse(output string, m *metricData) error {
	// Output example:
	// avg-cpu:  %user   %nice %system %iowait  %steal   %idle
	//           0.54    0.00    0.35    0.03    0.00   99.08
	//
	//Device:            tps   Blk_read/s   Blk_wrtn/s   Blk_read   Blk_wrtn
	//vda              11.55        83.04      1110.76    2390414   31975952
	var err error
	searchIndex := indexToSearch(output)
	tps, tpsIndex := findMetric(output, searchIndex)
	rSpeed, rIndex := findMetric(output, tpsIndex)
	wSpeed, _ := findMetric(output, rIndex)

	m.TPS, err = strconv.ParseFloat(tps, 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", tpsMetricName, "tps", metric.ErrParseFailed)
	}

	m.RSpeed, err = strconv.ParseFloat(rSpeed, 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", rSpeedMetricName, "r", metric.ErrParseFailed)
	}

	m.WSpeed, err = strconv.ParseFloat(wSpeed, 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", wSpeedMetricName, "w", metric.ErrParseFailed)
	}

	return nil
}

func findMetric(output string, firstIndex int) (string, int) {
	var startIndex, endIndex int
	notANumber := false
	for i := firstIndex; i < len(output); i++ {
		if unicode.IsDigit([]rune(output)[i]) {
			if !notANumber {
				startIndex = i
				notANumber = true
			}
			if !unicode.IsDigit([]rune(output)[i+1]) && !unicode.IsPunct([]rune(output)[i+1]) {
				endIndex = i
				break
			}
		}
	}

	return output[startIndex:endIndex], endIndex + 1
}

func indexToSearch(output string) int {
	blkIndex := strings.Index(output, "Blk_wrtn")
	var isNewWord bool
	for i := blkIndex; i < len(output); i++ {
		if unicode.IsSpace([]rune(output)[i]) {
			isNewWord = true
		}
		if isNewWord && !unicode.IsSpace([]rune(output)[i+1]) && unicode.IsDigit([]rune(output)[i+1]) {
			return i + 1
		}
	}
	return 0
}

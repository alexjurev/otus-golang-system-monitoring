package spaceinfo

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
	command = "df"
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
	//Filesystem           1K-blocks      Used Available Use% Mounted on
	//overlay               61202244  38931512  19129408  67% /
	//tmpfs                    65536         0     65536   0% /dev
	//shm                      65536         0     65536   0% /dev/shm
	///dev/vda1             61202244  38931512  19129408  67% /etc/resolv.conf
	///dev/vda1             61202244  38931512  19129408  67% /etc/hostname
	///dev/vda1             61202244  38931512  19129408  67% /etc/hosts
	var err error
	params, err := spaceInfoParams(output)
	if err != nil {
		return err
	}

	m.usedSpace, err = strconv.ParseFloat(strings.Trim(params[0], " "), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", usedSpaceMetricName, params[0], metric.ErrParseFailed)
	}
	m.percentSpace, err = strconv.ParseFloat(strings.Trim(params[1], " "), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", percentUsedMetricName, params[1], metric.ErrParseFailed)
	}
	m.usedISpace, err = strconv.ParseFloat(strings.Trim(params[2], " "), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", usedInodeSpaceMetricName, params[2], metric.ErrParseFailed)
	}
	m.percentIspace, err = strconv.ParseFloat(strings.Trim(params[3], " "), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", usedInodeSpaceMetricName, params[3], metric.ErrParseFailed)
	}

	return nil
}

func spaceInfoParams(output string) ([]string, error) {
	startIndex := strings.Index(output, "Mounted on")
	usedSpace, num := findNumber(output, startIndex)
	percentSpace, num := findNumber(output, num)
	usedISpace, num := findNumber(output, num)
	percentISpace, _ := findNumber(output, num)

	return []string{usedSpace, percentSpace, usedISpace, percentISpace}, nil
}

func findNumber(output string, startIndex int) (string, int) {
	var firstIndex int
	var isNewWord, firstFound bool
	for i := startIndex; i < len(output); i++ {
		if unicode.IsSpace([]rune(output)[i]) {
			isNewWord = true
		}
		if !firstFound && isNewWord && unicode.IsDigit([]rune(output)[i]) {
			firstIndex = i
			firstFound = true
		}
		if firstFound && isNewWord && unicode.IsSpace([]rune(output)[i+1]) {
			return output[firstIndex:i], i + 1
		}
	}
	return "", 0
}

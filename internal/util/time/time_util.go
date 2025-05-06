package time_util

import (
	"strings"
	"time"
)

var Now = time.Now

func ParseDurationToSeconds(input string) (int64, error) {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == 'h' || r == 'm' || r == 's'
	})

	units := []rune{'h', 'm', 's'}
	totalDuration := time.Duration(0)
	start := 0

	for i, part := range parts {
		if part == "" {
			continue
		}
		unit := string(units[i])
		durationStr := part + unit
		d, err := time.ParseDuration(durationStr)
		if err != nil {
			return 0, err
		}
		totalDuration += d
		start += len(part + unit)
	}

	return int64(totalDuration.Seconds()), nil
}

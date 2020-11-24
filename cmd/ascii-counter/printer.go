package main

import (
	"fmt"
	"os"

	"github.com/aybabtme/uniplot/histogram"

	"github.com/stainour/test12/internal/counter"
)

const width = 50

// printHistogram prints symbol frequency count as histogram into stdout.
func printHistogram(c counter.ASCIICodeCount) error {
	min, max, totalCount := c[0], c[0], 0
	var buckets = make([]histogram.Bucket, 0, len(c))

	for code, count := range c {
		if count == 0 {
			continue
		}

		totalCount += int(count)

		if count < min {
			min = count
		}

		if count > max {
			max = count
		}

		buckets = append(buckets, histogram.Bucket{
			Count: int(count),
			Min:   float64(code),
			Max:   float64(code),
		})
	}

	hist := histogram.Histogram{
		Min:     int(min),
		Max:     int(max),
		Count:   totalCount,
		Buckets: buckets,
	}

	return histogram.Fprintf(os.Stdout, hist, histogram.Linear(width), func(v float64) string {
		return fmt.Sprintf("%+q", string([]byte{byte(v)}))
	})
}

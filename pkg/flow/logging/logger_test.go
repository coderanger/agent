package logging_test

import (
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	gokitlevel "github.com/go-kit/log/level"
	"github.com/grafana/agent/pkg/flow/logging"
	flowlevel "github.com/grafana/agent/pkg/flow/logging/level"
	"github.com/stretchr/testify/require"
)

/* Most recent performance results on M2 Macbook Air:
$ go test -count=1 -benchmem ./pkg/flow/logging -run ^$ -bench BenchmarkLogging_
goos: darwin
goarch: arm64
pkg: github.com/grafana/agent/pkg/flow/logging
BenchmarkLogging_NoLevel_Prints-8             	  777288	      1533 ns/op	     392 B/op	      12 allocs/op
BenchmarkLogging_NoLevel_Drops-8              	25127053	        47.15 ns/op	      16 B/op	       1 allocs/op
BenchmarkLogging_GoKitLevel_Drops_Sprintf-8   	 3216174	       373.4 ns/op	     352 B/op	       9 allocs/op
BenchmarkLogging_GoKitLevel_Drops-8           	 5947682	       201.0 ns/op	     480 B/op	       6 allocs/op
BenchmarkLogging_GoKitLevel_Prints-8          	  683368	      1746 ns/op	     873 B/op	      17 allocs/op
BenchmarkLogging_Slog_Drops-8                 	32123890	        36.71 ns/op	      16 B/op	       1 allocs/op
BenchmarkLogging_Slog_Prints-8                	 1000000	      1173 ns/op	      47 B/op	       3 allocs/op
BenchmarkLogging_FlowLevel_Drops-8            	16338895	        74.42 ns/op	     176 B/op	       3 allocs/op
BenchmarkLogging_FlowLevel_Prints-8           	  682309	      1740 ns/op	     873 B/op	      17 allocs/op
*/

const testStr = "this is a test string"

func BenchmarkLogging_NoLevel_Prints(b *testing.B) {
	logger, err := logging.New(io.Discard, debugLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		logger.Log("msg", "test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func BenchmarkLogging_NoLevel_Drops(b *testing.B) {
	logger, err := logging.New(io.Discard, warnLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		logger.Log("msg", "test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func BenchmarkLogging_GoKitLevel_Drops_Sprintf(b *testing.B) {
	logger, err := logging.New(io.Discard, debugLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		gokitlevel.Debug(logger).Log("msg", fmt.Sprintf("test message %d, error=%v, str=%s, duration=%v", i, testErr, testStr, time.Since(start)))
	}
}

func BenchmarkLogging_GoKitLevel_Drops(b *testing.B) {
	logger, err := logging.New(io.Discard, debugLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		gokitlevel.Debug(logger).Log("msg", "test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func BenchmarkLogging_GoKitLevel_Prints(b *testing.B) {
	logger, err := logging.New(io.Discard, debugLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	testStr := "this is a test string"
	start := time.Now()
	for i := 0; i < b.N; i++ {
		gokitlevel.Warn(logger).Log("msg", "test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func BenchmarkLogging_Slog_Drops(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	testErr := fmt.Errorf("test error")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		logger.Debug("test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func BenchmarkLogging_Slog_Prints(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	testErr := fmt.Errorf("test error")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func BenchmarkLogging_FlowLevel_Drops(b *testing.B) {
	logger, err := logging.New(io.Discard, debugLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		flowlevel.Debug(logger).Log("msg", "test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func BenchmarkLogging_FlowLevel_Prints(b *testing.B) {
	logger, err := logging.New(io.Discard, debugLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		flowlevel.Info(logger).Log("msg", "test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func debugLevelOptions() logging.Options {
	opts := logging.Options{}
	opts.SetToDefault()
	opts.Level = logging.LevelInfo
	return opts
}

func warnLevelOptions() logging.Options {
	opts := debugLevelOptions()
	opts.Level = logging.LevelWarn
	return opts
}

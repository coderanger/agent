package logging_test

import (
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	gokitlevel "github.com/go-kit/log/level"
	"github.com/grafana/agent/pkg/flow/logging"
	"github.com/stretchr/testify/require"
)

/* Most recent performance results on M2 Macbook Air:
$ go test -count=1 -benchmem ./pkg/flow/logging -run ^$ -bench BenchmarkLogging_
goos: darwin
goarch: arm64
pkg: github.com/grafana/agent/pkg/flow/logging
BenchmarkLogging_NoLevel_Prints-8                 765578              1958 ns/op             408 B/op         13 allocs/op
BenchmarkLogging_NoLevel_Drops-8                 3273050               335.6 ns/op            64 B/op          7 allocs/op
BenchmarkLogging_GoKitLevel_Drops_Sprintf-8      2479004               486.1 ns/op           464 B/op         11 allocs/op
BenchmarkLogging_GoKitLevel_Drops-8              2389933               504.1 ns/op           528 B/op         11 allocs/op
BenchmarkLogging_GoKitLevel_Prints-8              680679              1739 ns/op             873 B/op         17 allocs/op
BenchmarkLogging_Slog_Drops-8                   23788833                50.28 ns/op           32 B/op          2 allocs/op
BenchmarkLogging_Slog_Prints-8                   1000000              1195 ns/op              64 B/op          5 allocs/op
*/

func BenchmarkLogging_NoLevel_Prints(b *testing.B) {
	logger, err := logging.New(io.Discard, debugLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	testStr := "this is a test string"
	start := time.Now()
	for i := 0; i < b.N; i++ {
		logger.Log("msg", "test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func BenchmarkLogging_NoLevel_Drops(b *testing.B) {
	logger, err := logging.New(io.Discard, warnLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	testStr := "this is a test string"
	start := time.Now()
	for i := 0; i < b.N; i++ {
		logger.Log("msg", "test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
	}
}

func BenchmarkLogging_GoKitLevel_Drops_Sprintf(b *testing.B) {
	logger, err := logging.New(io.Discard, debugLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	testStr := "this is a test string"
	start := time.Now()
	for i := 0; i < b.N; i++ {
		gokitlevel.Debug(logger).Log("msg", fmt.Sprintf("test message %d, error=%v, str=%s, duration=%v", i, testErr, testStr, time.Since(start)))
	}
}

func BenchmarkLogging_GoKitLevel_Drops(b *testing.B) {
	logger, err := logging.New(io.Discard, debugLevelOptions())
	require.NoError(b, err)

	testErr := fmt.Errorf("test error")
	testStr := "this is a test string"
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
	testStr := "this is a test string"
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
	testStr := "this is a test string"
	start := time.Now()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "i", i, "err", testErr, "str", testStr, "duration", time.Since(start))
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

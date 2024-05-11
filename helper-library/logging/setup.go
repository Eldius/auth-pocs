package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
)

var (
	logKeys = []string{
		"host",
		"service.name",
		"level",
		"message",
		"time",
		"error",
		"source",
		"function",
		"file",
		"line",
	}
)

func SetupLogs(serviceName string, debug bool) error {
	var h slog.Handler
	var w io.Writer = os.Stdout

	replaceAttrFunc := func(groups []string, a slog.Attr) slog.Attr {
		if slices.Contains(logKeys, a.Key) {
			return a
		}
		if strings.HasPrefix(a.Key, "request.") || strings.HasPrefix(a.Key, "response.") {
			return a
		}
		if a.Key == "msg" {
			a.Key = "message"
			return a
		}
		a.Key = fmt.Sprintf("custom.%s.%s", serviceName, a.Key)
		return a
	}

	h = slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   true,
		Level:       parseLogLevel(debug),
		ReplaceAttr: replaceAttrFunc,
	})
	logger := slog.New(h)
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	slog.SetDefault(logger.With(
		slog.String("service.name", serviceName),
		slog.String("host", host),
	))

	return nil
}

func parseLogLevel(debug bool) slog.Level {
	if debug {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

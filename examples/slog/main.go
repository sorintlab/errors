package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/phsym/console-slog"
	"github.com/sorintlab/errors"
)

var detailedErrors = true

func Error(err error) slog.Attr {
	if detailedErrors {
		return slog.Group("error",
			slog.String("message", err.Error()),
			slog.String("details", "\n"+strings.Join(errors.PrintErrorDetails(err), "\n")),
		)
	}

	return slog.Group("error",
		slog.String("message", err.Error()),
	)
}

func main() {
	logger := slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelDebug}))

	err1 := errors.New("error 1")
	err2 := errors.Wrap(err1, "error 2 wrapping error 1")
	err3 := errors.New("error 3")
	err4 := errors.Join(err2, err3)
	err5 := errors.New("error 5")
	err6 := errors.Join(err4, err5)
	err7 := errors.Wrap(err6, "error 7")
	err8 := errors.WithStack(err7)

	logger.Error("error", Error(err8))
}

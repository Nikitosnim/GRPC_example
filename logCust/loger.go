package logcust

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type DesignHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type DesignHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *DesignHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	// Adding colors to indicate levels
	switch r.Level {
	case slog.LevelDebug:
		level = color.CyanString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	// Processing additional arguments
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	// make arguments in the form json
	b, err := json.MarshalIndent(fields, "", "   ")
	if err != nil {
		return err
	}

	// time formatting
	time := r.Time.Format("[Jan_2 15:04:05]")

	// conclusion
	h.l.Println(time, level, r.Message, color.WhiteString(string(b)))

	return nil
}

// Constructor for handlers format text
func NewTextDesignHandler(
	out io.Writer,
	opts DesignHandlerOptions,
) *DesignHandler {
	h := &DesignHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

// Constructor for handlers format JSON
func NewJSONDesignHandler(
	out io.Writer,
	opts DesignHandlerOptions,
) *DesignHandler {
	h := &DesignHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

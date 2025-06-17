package logging

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	types "github.com/pocketbase/pocketbase/tools/types"
)

type HookHandlerOptions struct {
	Level     slog.Leveler
	WriteFunc func(ctx context.Context, record *Record) error
}

var _ slog.Handler = (*HookHandler)(nil)

type HookHandler struct {
	mutex   *sync.Mutex
	parent  *HookHandler
	options *HookHandlerOptions
	group   string
	attrs   []slog.Attr
}

func NewHookHandler(opts *HookHandlerOptions) *HookHandler {
	if opts == nil {
		opts = &HookHandlerOptions{}
	}

	h := &HookHandler{
		mutex:   &sync.Mutex{},
		options: opts,
	}

	if h.options.WriteFunc == nil {
		panic("`options.WriteFunc` is nil")
	}

	if h.options.Level == nil {
		h.options.Level = slog.LevelInfo
	}

	return h
}

func (h *HookHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.options.Level.Level()
}

func (h *HookHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	return &HookHandler{
		parent:  h,
		mutex:   h.mutex,
		options: h.options,
		group:   name,
	}
}

func (h *HookHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	return &HookHandler{
		parent:  h,
		mutex:   h.mutex,
		options: h.options,
		attrs:   attrs,
	}
}

func (h *HookHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.group != "" {
		h.mutex.Lock()
		attrs := make([]any, 0, len(h.attrs)+r.NumAttrs())
		for _, a := range h.attrs {
			attrs = append(attrs, a)
		}
		h.mutex.Unlock()

		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, a)
			return true
		})

		r = slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
		r.AddAttrs(slog.Group(h.group, attrs...))
	} else if len(h.attrs) > 0 {
		r = r.Clone()

		h.mutex.Lock()
		r.AddAttrs(h.attrs...)
		h.mutex.Unlock()
	}

	if h.parent != nil {
		return h.parent.Handle(ctx, r)
	}

	data := make(map[string]any, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		if err := h.resolveAttr(data, a); err != nil {
			return false
		}
		return true
	})

	log := &Record{
		Time:    r.Time,
		Message: r.Message,
		Data:    types.JSONMap[any](data),
	}
	switch r.Level {
	case slog.LevelDebug:
		log.Level = LevelDebug
	case slog.LevelInfo:
		log.Level = LevelInfo
	case slog.LevelWarn:
		log.Level = LevelWarn
	case slog.LevelError:
		log.Level = LevelError
	default:
		log.Level = Level(fmt.Sprintf("LV(%d)", r.Level))
	}

	if err := h.writeRecord(ctx, log); err != nil {
		return err
	}

	return nil
}

func (h *HookHandler) SetLevel(level slog.Level) {
	h.mutex.Lock()
	h.options.Level = level
	h.mutex.Unlock()
}

func (h *HookHandler) writeRecord(ctx context.Context, r *Record) error {
	if h.parent != nil {
		return h.parent.writeRecord(ctx, r)
	}

	return h.options.WriteFunc(ctx, r)
}

func (h *HookHandler) resolveAttr(data map[string]any, attr slog.Attr) error {
	attr.Value = attr.Value.Resolve()

	if attr.Equal(slog.Attr{}) {
		return nil
	}

	switch attr.Value.Kind() {
	case slog.KindGroup:
		{
			attrs := attr.Value.Group()
			if len(attrs) == 0 {
				return nil
			}

			groupData := make(map[string]any, len(attrs))

			for _, subAttr := range attrs {
				h.resolveAttr(groupData, subAttr)
			}

			if len(groupData) > 0 {
				data[attr.Key] = groupData
			}
		}

	default:
		{
			switch v := attr.Value.Any().(type) {
			case error:
				data[attr.Key] = v.Error()
			default:
				data[attr.Key] = v
			}
		}
	}

	return nil
}

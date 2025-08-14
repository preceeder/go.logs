package logs

import (
	"context"
	"github.com/preceeder/go.base"
	"log/slog"
)

type MoreHandler struct {
	//TransparentTransmission bool // 日志是否往高等级传递 LevelErr(8) -> LevelWarn(4) -> LevelInfo(0) -> LevelDebug(-4)
	MHandler []Config   // slog.Leve   从slog.Level 低到高排序
	MinLevel slog.Level // 最小的 level    如果有 debug, info 最小的就是 info  	LevelDebug = -4 ,LevelInfo = 0 ,LevelWarn = 4 ,LevelError = 8
}

func NewMoreHandler(w []Config, opts *slog.HandlerOptions) *MoreHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	minLevel := w[len(w)-1].LogLevel
	handler := MoreHandler{MinLevel: minLevel, MHandler: []Config{}}
	for _, wl := range w {
		if wl.OutType == "json" || wl.OutType == "" {
			wl.l = slog.NewJSONHandler(wl.w, opts)
		} else {
			wl.l = slog.NewTextHandler(wl.w, opts)
		}
		handler.MHandler = append(handler.MHandler, wl)
	}
	return &handler
}

func (h *MoreHandler) Enabled(_ context.Context, level slog.Level) bool {
	if level >= h.MinLevel {
		return true
	}
	return false
}

func (h *MoreHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MoreHandler{}
}

func (h *MoreHandler) WithGroup(name string) slog.Handler {
	return &MoreHandler{}
}

func (h *MoreHandler) Handle(c context.Context, r slog.Record) error {
	if vc, ok := c.(base.BaseContext); ok {
		r.Add("REQUESTID", vc.GetRequestId())
		r.Add("USERID", vc.GetUserId())
	}
	for _, handler := range h.MHandler {
		if handler.LogLevel > r.Level {
			continue
		}

		if handler.LogLevel == r.Level {
			_ = handler.l.Handle(c, r)
			if handler.TransparentTransmission {
				continue
			}
			break
		}
		_ = handler.l.Handle(c, r)
		if !handler.TransparentTransmission {
			break
		}
	}

	return nil
}

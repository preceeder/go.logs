package logs

import (
	"github.com/preceeder/go.logs/lumberjack"
	"io"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Config struct {
	FileName                string       `json:"fileName"`
	LogLevel                slog.Level   `json:"logLevel"`                // 	LevelDebug = -4 |  LevelInfo  = 0 | LevelWarn  = 4 | LevelError  = 8
	MaxSize                 int          `json:"maxSize"`                 // 日志轮转的最大size 单位： MB,  -1 不限制大小， 0 默认100MB,
	RotateTime              string       `json:"rotateTime"`              // 日志轮转 每天的特定时间点  12:00:00, MaxSize 有值的时候达到MaxSize限制也会轮转
	MaxAge                  int          `json:"maxAge"`                  // 历史日志保存的最大时间 天
	MaxBackups              int          `json:"maxBackups"`              // 保存的 历史日志数量
	TransparentTransmission bool         `json:"transparentTransmission"` // 日志输出是否传递到高等级等级  level 越小等级越高, ErrorFileName 日志可以写入到 WarnFileName|InfoFileName|DebugFileName
	StdOut                  bool         `json:"stdOut"`                  // 是否在终端输出   0 不输出, 1 输出
	Compress                bool         `json:"compress"`                // 备份的日志是否压缩
	OutType                 string       `json:"outType"`                 // json, txt    输出格式， 默认json
	l                       slog.Handler // 内部填充
	w                       io.Writer
}
type SlogConfig struct {
	Config           []Config `json:"config"`
	NotSetDefaultLog bool     `json:"notSetDefaultLog"` // 不设置为默认的slog   true 不设置， false设置
}

var DefaultSlogConfig SlogConfig = SlogConfig{
	Config: []Config{
		{
			FileName: "",
			LogLevel: slog.LevelInfo,
			StdOut:   true,
		},
	},
	NotSetDefaultLog: false,
}

func NewSlog(cfg *SlogConfig) *slog.Logger {
	return mHandlerSlog(cfg)
}

func mHandlerSlog(cfg *SlogConfig) *slog.Logger {
	if cfg == nil {
		cfg = &DefaultSlogConfig
	}
	iohandler := getIoWriter(cfg)

	opt := &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			d := a.Value.Any().(*slog.Source)
			a.Value = slog.AnyValue(strings.Join([]string{d.File, strconv.Itoa(d.Line)}, ":"))
		} else if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format("20060102150405.000000"))
		}
		return a
	}}

	var log *slog.Logger
	log = slog.New(NewMoreHandler(iohandler, opt))
	if !cfg.NotSetDefaultLog {
		slog.SetDefault(log)
	}
	return log
}

func getIoWriter(cfg *SlogConfig) []Config {
	iohandler := []Config{}
	for _, cf := range cfg.Config {
		if cf.FileName != "" {
			el := &lumberjack.Logger{
				Filename:   cf.FileName,
				MaxSize:    cf.MaxSize,
				MaxBackups: cf.MaxBackups,
				MaxAge:     cf.MaxAge,
				Compress:   cf.Compress,
				RotateTime: cf.RotateTime,
				LocalTime:  true,
			}
			el.Init()
			cf.w = el
		}

		if cf.StdOut {
			if cf.w != nil {
				cf.w = io.MultiWriter(cf.w, os.Stdout)
			} else {
				cf.w = os.Stdout
			}
		}
		iohandler = append(iohandler, cf)
	}
	slices.SortFunc(iohandler, func(a, b Config) int {
		if a.LogLevel < b.LogLevel {
			return -1
		}
		if a.LogLevel > b.LogLevel {
			return 1
		}
		return 0
	})
	return iohandler
}

// json 字符串可以不加转议符"\" 输出
type LogStr string

func (d LogStr) MarshalJSON() ([]byte, error) {
	return []byte(d), nil
}

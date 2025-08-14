package logs

import (
	"encoding/json"
	"fmt"
	"github.com/preceeder/go.base"
	"log/slog"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Dar string

func (d Dar) MarshalJSON() ([]byte, error) {
	return []byte(d), nil
}

func TestLou(t *testing.T) {
	tm := time.NewTimer(time.Second * 5)
	for {
		select {
		case <-tm.C:
			fmt.Println("timed on")
		}
		fmt.Println(tm.Reset(time.Second * 5))
	}
}

func TestLogs(t *testing.T) {
	cfg := SlogConfig{
		Config: []Config{
			{
				FileName:                "logs/error.log",
				LogLevel:                slog.LevelError,
				TransparentTransmission: true,
				StdOut:                  true,
				MaxSize:                 1,
				MaxAge:                  3,
				MaxBackups:              4,
				RotateTime:              "15:19:00",
			},
			{
				FileName:                "logs/info.log",
				LogLevel:                slog.LevelInfo,
				OutType:                 "txt",
				TransparentTransmission: true,
				StdOut:                  false,
				MaxSize:                 1,
				MaxAge:                  3,
				MaxBackups:              4,
				RotateTime:              "15:22:00",
			},
			{
				FileName:                "logs/warn.log",
				LogLevel:                slog.LevelWarn,
				TransparentTransmission: false,
				StdOut:                  true,
				MaxSize:                 1,
				MaxAge:                  3,
				MaxBackups:              4,
				RotateTime:              "15:24:00",
			},
		},
		NotSetDefaultLog: false,
	}
	NewSlog(&cfg)
	//i := 1000
	//j := 0
	ctx := base.Context{RequestId: "sdsd"}
	slog.ErrorContext(ctx, "sds")
	//time.Sleep(time.Second * 100)
	marshal, _ := json.Marshal(map[string]any{"Name": "ois", "Age": 23})
	de := Dar(marshal)
	slog.InfoContext(ctx, "", "data", de)
}

func BenchmarkTestNewSlog(b *testing.B) {
	cfg := SlogConfig{
		Config: []Config{
			{
				FileName:                "logs/error.log",
				LogLevel:                slog.LevelError,
				TransparentTransmission: true,
				StdOut:                  true,
				MaxSize:                 1,
				MaxAge:                  3,
				MaxBackups:              4,
			},
			{
				FileName:                "logs/info.log",
				LogLevel:                slog.LevelInfo,
				TransparentTransmission: true,
				StdOut:                  false,
				MaxSize:                 1,
				MaxAge:                  3,
				MaxBackups:              4,
			},
			{
				FileName:                "logs/warn.log",
				LogLevel:                slog.LevelWarn,
				TransparentTransmission: false,
				StdOut:                  true,
				MaxSize:                 1,
				MaxAge:                  3,
				MaxBackups:              4,
			},
		},
		NotSetDefaultLog: false,
	}
	NewSlog(&cfg)
	//i := 1000
	//j := 0
	ctx := base.Context{RequestId: "sdsd"}
	slog.InfoContext(ctx, "sds")

	marshal, _ := json.Marshal(map[string]any{"Name": "ois", "Age": 23})
	de := Dar(marshal)
	slog.InfoContext(ctx, "", "data", de)
	//for j <= b.N {
	//	j++
	//	slog.InfoContext(ctx, "sds")
	//	slog.ErrorContext(ctx, "hahah", "error", "ssss")
	//	//slog.Info("", "sss", j)
	//}

}

func BenchmarkNewSlog(b *testing.B) {
	file := "scabhvuiehviuusdvuauisgvuiav"
	line := 3
	//j := 0
	//for j <= b.N {
	//	j++
	//	//fmt.Sprintf("%s:%d", file, line)
	//	//
	//	_ = strings.Join([]string{file, strconv.Itoa(line)}, ":")
	//
	//}

	dd := strings.Join([]string{file, strconv.Itoa(line)}, ":")
	fmt.Println(dd)
}

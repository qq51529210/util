package util

import "github.com/qq51529210/log"

var (
	logFile *log.File
)

// LogCfg 日志的配置
type LogCfg struct {
	// 日志保存的根目录
	RootDir string `json:"rootDir" yaml:"rootDir" validate:"required,filepath"`
	// 每一份日志文件的最大字节，使用 1.5/K/M/G/T 这样的字符表示。
	MaxFileSize string `json:"maxFileSize" yaml:"maxFileSize"`
	// 保存的最大天数，最小是1天。
	MaxKeepDay float64 `json:"maxKeepDay" yaml:"maxKeepDay"`
	// 同步到磁盘的时间间隔，单位，毫秒。最小是10毫秒。
	SyncInterval int `json:"syncInterval" yaml:"syncInterval" validate:"required,min=1"`
	// 是否输出到控制台，out/err
	Std string `json:"std" yaml:"std" validate:"omitempty,oneof=out err"`
	// 禁用的日志级别
	DisableLevel []string `json:"disableLevel" yaml:"disableLevel" validate:"omitempty,dive,oneof=debug info warn error"`
}

func InitLog(cfg *LogCfg) error {
	var logCfg log.FileConfig
	CopyStruct(&logCfg, cfg)
	f, err := log.NewFile(&logCfg)
	if err != nil {
		return err
	}
	logger := log.GetLogger()
	for i := 0; i < len(cfg.DisableLevel); i++ {
		if cfg.DisableLevel[i] == "debug" {
			logger.EnableDebug(false)
			continue
		}
		if cfg.DisableLevel[i] == "info" {
			logger.EnableDebug(false)
			continue
		}
		if cfg.DisableLevel[i] == "warn" {
			logger.EnableDebug(false)
			continue
		}
		if cfg.DisableLevel[i] == "error" {
			logger.EnableDebug(false)
			continue
		}
	}
	logFile = f
	//
	return nil
}

func CloseLog() {
	if logFile != nil {
		logFile.Close()
	}
}

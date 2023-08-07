package util

import (
	"io"

	"github.com/qq51529210/log"
)

// LogCfg 日志的配置
type LogCfg struct {
	log.FileConfig `yaml:",inline"`
	// 禁用的日志级别
	DisableLevel []string `json:"disableLevel" yaml:"disableLevel" validate:"omitempty,dive,oneof=debug info warn error"`
}

// NewFileLog 返回日志文件
func NewFileLog(cfg *LogCfg) (io.WriteCloser, error) {
	file, err := log.NewFile(&cfg.FileConfig)
	if err != nil {
		return nil, err
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
	//
	return file, nil
}

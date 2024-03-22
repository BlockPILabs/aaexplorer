package task

import (
	interConfig "github.com/BlockPILabs/aaexplorer/config"
	interlog "github.com/BlockPILabs/aaexplorer/internal/log"
)

var logger = interlog.L()

func SetLogger(lg interlog.Logger) {
	logger = lg
}

var config = interConfig.DefaultConfig()

func SetConfig(cfg *interConfig.Config) {
	config = cfg
}

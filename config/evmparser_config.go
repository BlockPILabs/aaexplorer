package config

import (
	_ "embed"
)

//go:embed abi/erc4337-abi.json
var ERC_4337_ABI string

//go:embed abi/erc4337-abi-v0.7.json
var ERC_4337_ABI_V07 string

var HandleOpsMap = map[string]string{
	"0x1fad948c": "0.6",
	"0x765e827f": "0.7",
}

type EvmParserConfig struct {
	StartBlock map[string]int64 `mapstructure:"startBlock" toml:"startBlock"` // -1 start by latest , 0 start by first , >0 start by set
	Multi      int              `mapstructure:"multi" toml:"multi"`
	Batch      int              `mapstructure:"batch" toml:"batch"`
}

func DefaultEvmParserConfig() *EvmParserConfig {
	return &EvmParserConfig{
		StartBlock: map[string]int64{},
		Multi:      10,
		Batch:      10,
	}
}

func (c *EvmParserConfig) GetAbi(version string) string {
	if version == "0.6" {
		return ERC_4337_ABI
	} else if version == "0.7" {
		return ERC_4337_ABI_V07
	}
	return ""
}

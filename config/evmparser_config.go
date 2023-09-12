package config

import (
	_ "embed"
)

//go:embed abi/erc4337-abi.json
var ERC_4337_ABI string

type EvmParserConfig struct {
	StartBlock map[string]int64 `mapstructure:"startBlock" toml:"startBlock"` // -1 start by latest , 0 start by first , >0 start by set
	Multi      int              `mapstructure:"multi" toml:"multi"`
	Batch      int              `mapstructure:"batch" toml:"batch"`
	Abi        string           `mapstructure:"abi" toml:"abi"`
}

func DefaultEvmParserConfig() *EvmParserConfig {
	return &EvmParserConfig{
		StartBlock: map[string]int64{},
		Multi:      10,
		Batch:      10,
		Abi:        ERC_4337_ABI,
	}
}

func (c *EvmParserConfig) GetAbi() string {
	if len(c.Abi) < 1 {
		return ERC_4337_ABI
	}
	return c.Abi
}

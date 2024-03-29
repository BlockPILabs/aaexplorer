package config

import "math"

const (
	Default = "default"

	PolygonMumbai = "polygon-mumbai"
	Polygon       = "polygon"
	Eth           = "eth"
	BSC           = "bsc"

	EthNative      = "eth"
	EvmDecimal     = 18
	BscNative      = "bnb"
	PolygonNative  = "matic"
	EthNetwork     = "ethereum"
	BscNetwork     = "bsc"
	PolygonNetwork = "polygon"

	WETH   = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	WMATIC = "0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"
	WBNB   = "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"

	WhaleUsd           = 50000
	AddressTypeAccount = 1
	AddressTypeToken   = 2
	ZeroAddress        = "0x0000000000000000000000000000000000000000"

	RangeH24 = "h24"
	RangeD7  = "d7"
	RangeD30 = "d30"

	AssetExpireTime = 2 * 3600 * 1000
	AnalyzeTop7     = 7
	FeePrecision    = 6

	TopNum = 300
)

const (
	MaxPage = math.MaxInt32
	MinPage = 1

	DefaultPerPage  = 10
	MaxPerPage      = 1000
	CreateInBatches = 1000
	MinPerPage      = 0

	OrderAsc  = 1
	OrderDesc = -1
)

const TokenTypeBase = "base"

const (
	ModeLocal   = "local"
	ModeDevnet  = "devnet"
	ModeTestnet = "testnet"
	ModeMainnet = "mainnet"
)

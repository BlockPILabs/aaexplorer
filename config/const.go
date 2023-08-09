package config

import "math"

const (
	Default = "default"

	PolygonMumbai = "polygon-mumbai"
	Polygon       = "polygon"
	Eth           = "eth"
	BSC           = "bsc"
)

const (
	MaxPage = math.MaxInt32
	MinPage = 1

	MaxPerPage      = 1000
	CreateInBatches = 1000
	MinPerPage      = 0

	OrderAsc  = 1
	OrderDesc = -1
)

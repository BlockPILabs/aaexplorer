package vo

import "github.com/shopspring/decimal"

type WhaleOverviewRequest struct {
	Network string `json:"network"`
}

type WhaleOverviewResponse struct {
	TxDominance   decimal.Decimal `json:"txDominance"`
	Ratio         decimal.Decimal `json:"ratio"`
	TotalAssetUsd decimal.Decimal `json:"totalAssetUsd"`
	TotalEthUsd   decimal.Decimal `json:"totalEthUsd"`
}

type WhaleChartRequest struct {
	Network   string `json:"network"`
	TimeRange string `json:"timeRange"`
}

type WhaleChartResponse struct {
	Details []*WhaleChartDetail `json:"details"`
}

type WhaleChartDetail struct {
	Time  int64           `json:"time"`
	Value decimal.Decimal `json:"value"`
}

type ByWhaleDailyStatisticTime []*WhaleChartDetail

func (b ByWhaleDailyStatisticTime) Len() int           { return len(b) }
func (b ByWhaleDailyStatisticTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByWhaleDailyStatisticTime) Less(i, j int) bool { return b[i].Time < b[j].Time }

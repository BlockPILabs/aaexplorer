package schema

import (
	"entgo.io/ent"
	"time"
)

type HotAATokenStatistic struct {
	Id            int64
	TokenSymbol   string
	Network       string
	StatisticType string
	TotalNum      int64
	CreateTime    time.Time
	ent.Schema
}

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"time"
)

type DailyStatisticDay struct {
	Id            int64
	Network       string
	StatisticTime time.Time
	TxNum         int64
	UserOpsNum    int64
	GasFee        decimal.Decimal
	ActiveWallet  int64
	CreateTime    time.Time
	ent.Schema
}

func (DailyStatisticDay) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Time("statistic_time").
			StructTag(`json:"statisticTime"`),
		field.Int64("tx_num").
			StructTag(`json:"txNum"`),
		field.Int64("user_ops_num").
			StructTag(`json:"userOpsNum"`),
		field.Int64("gas_fee").
			StructTag(`json:"gasFee"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("active_wallet").
			StructTag(`json:"activeWallet"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (DailyStatisticDay) Edges() []ent.Edge {
	return nil
}

func (DailyStatisticDay) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "daily_statistic_day"},
	}
}

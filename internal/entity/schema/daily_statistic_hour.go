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

type DailyStatisticHour struct {
	Id            int64
	Network       string
	StatisticTime int64
	TxNum         int64
	UserOpsNum    int64
	GasFee        decimal.Decimal
	ActiveWallet  int64
	CreateTime    time.Time
	ent.Schema
}

func (DailyStatisticHour) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Int64("statistic_time").
			StructTag(`json:"statisticTime"`),
		field.Int64("tx_num").
			StructTag(`json:"txNum"`),
		field.Int64("user_ops_num").
			StructTag(`json:"userOpsNum"`),
		field.Int64("gas_fee").
			StructTag(`json:"gasFee"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("gas_fee_usd").
			StructTag(`json:"gasFeeUsd"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("active_wallet").
			StructTag(`json:"activeWallet"`),
		field.Int64("paymaster_gas_paid").
			StructTag(`json:"paymasterGasPaid"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("paymaster_gas_paid_usd").
			StructTag(`json:"paymasterGasPaidUsd"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("bundler_gas_profit").
			StructTag(`json:"bundlerGasProfit"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("bundler_gas_profit_usd").
			StructTag(`json:"bundlerGasProfitUsd"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (DailyStatisticHour) Edges() []ent.Edge {
	return nil
}

func (DailyStatisticHour) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "daily_statistic_hour"},
	}
}

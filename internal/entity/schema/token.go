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

type Token struct {
	ent.Schema
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("contract_address").
			MaxLen(255).
			StructTag(`json:"contractAddress"`),
		field.String("symbol").
			MaxLen(255).
			StructTag(`json:"symbol"`),
		field.String("full_name").
			MaxLen(255).
			StructTag(`json:"fullName"`),
		field.Int64("token_price").StructTag(`json:"tokenPrice"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("last_time").
			StructTag(`json:"lastTime"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updateTime"`).
			Immutable(),
		field.Int64("market_rank").
			StructTag(`json:"marketRank"`),
		field.String("type").
			StructTag(`json:"type"`).
			Immutable().Nillable(),
	}
}

func (Token) Edges() []ent.Edge {
	return nil
}

func (Token) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "token"},
	}
}

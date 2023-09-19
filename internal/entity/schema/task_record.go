package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type TaskRecord struct {
	ID         int64     `db:"id"`
	Network    string    `db:"network"`
	TaskType   string    `db:"task_type"`
	LastTime   time.Time `db:"last_time"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
	ent.Schema
}

func (TaskRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("task_type").
			StructTag(`json:"taskType"`),
		field.Time("last_time").
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
	}
}

func (TaskRecord) Edges() []ent.Edge {
	return nil
}

func (TaskRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "task_record"},
	}
}

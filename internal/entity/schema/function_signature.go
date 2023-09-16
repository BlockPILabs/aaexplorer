package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type FunctionSignature struct {
	ent.Schema
}

/*
create table public.function_signature
(

	signature   varchar(16) not null
	    constraint function_signature_pk
	        primary key,
	name        varchar(255),
	text        text,
	bytes       bytea,
	create_time timestamp with time zone

);

alter table public.function_signature

	owner to postgres;
*/
func (FunctionSignature) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("signature").StructTag(`json:"signature"`),
		field.String("name").StructTag(`json:"name"`).Optional(),
		field.String("text").StructTag(`json:"text"`).Optional(),
		field.Bytes("bytes").StructTag(`json:"bytes"`).Optional(),
		field.Time("create_time").StructTag(`json:"createTime"`).Optional(),
	}
}

func (FunctionSignature) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "function_signature"},
	}
}

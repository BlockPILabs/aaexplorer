package entity

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --feature sql/schemaconfig --feature sql/lock --feature sql/execquery --feature sql/modifier --target ./ent ./schema

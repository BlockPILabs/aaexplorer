package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
)

type baseDao struct {
}

func (dao *baseDao) sortField(ctx context.Context, fields []string, sort int) string {
	if sort >= 0 && len(fields) > sort {
		return fields[sort]
	}
	return ""
}

func (dao *baseDao) orderOptions(ctx context.Context, order int) (opts []sql.OrderTermOption) {
	if order > 0 {
		opts = append(opts, sql.OrderAsc(), sql.OrderNullsFirst())

	} else if order < 0 {
		opts = append(opts, sql.OrderDesc(), sql.OrderNullsLast())
	}
	return nil
}

package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type aaAccountDao struct {
	baseDao
}

var AaAccountDao = &aaAccountDao{}

func (dao *aaAccountDao) Search(ctx context.Context, tx *ent.Client, req vo.SearchAllRequest) (a ent.AaAccountDataSlice, err error) {

	types := []string{
		config.AaAccountTypeAA,
		config.AaAccountTypeFactory,
		config.AaAccountTypePaymaster,
		config.AaAccountTypeBundler,
		//config.AaAccountTypeEntryPoint,
		//"",
	}
	term := req.Term
	if !utils.Has0xPrefix(term) {
		term = "0x" + term
	}

	q := func(p *sql.Predicate) *sql.Predicate {
		if utils.IsHexAddress(term) {
			p.EQ(aaaccountdata.FieldID, term)
		} else {
			p.HasSuffix(aaaccountdata.FieldID, term)
		}
		return p
	}

	sel := sql.Select("*").From(sql.Table(aaaccountdata.Table)).Where(
		q(sql.P().EQ(aaaccountdata.FieldAaType, "")),
	).Limit(10)

	for _, s := range types {
		sel = sel.UnionAll(
			sql.Select("*").From(sql.Table(aaaccountdata.Table)).Where(
				q(sql.P().EQ(aaaccountdata.FieldAaType, s)),
			).Limit(10),
		)
	}

	query, vars := sel.Query()
	rows, err := tx.QueryContext(ctx, query, vars...)
	if err != nil {
		return nil, err
	}
	err = rows.Scan(a)
	if err != nil {
		return a, err
	}

	return a, err
}

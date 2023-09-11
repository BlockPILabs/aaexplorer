package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"strings"
)

type aaAccountDao struct {
	baseDao
}

var AaAccountDao = &aaAccountDao{}

func (dao *aaAccountDao) Search(ctx context.Context, tx *ent.Client, req vo.SearchAllRequest) (a ent.AaAccountDataSlice, err error) {

	term := strings.ToLower(req.Term)
	if !utils.Has0xPrefix(term) {
		term = "0x" + term
	}

	query := tx.AaAccountData.Query()
	if utils.IsHexAddress(term) {
		query.Where(
			aaaccountdata.IDEQ(term),
		)
	} else {
		query.Where(
			sql.FieldHasPrefix(aaaccountdata.FieldID, term),
		)
	}

	return query.Limit(100).All(ctx)
}

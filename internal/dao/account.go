package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
)

type accountDao struct {
	baseDao
}

var AccountDao = &accountDao{}

func (dao *accountDao) GetAbiByAddress(ctx context.Context, tx *ent.Client, address string) (account *ent.Account, err error) {
	return tx.Account.Get(ctx, address)
}

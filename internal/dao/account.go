package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/account"
)

type accountDao struct {
	baseDao
}

var AccountDao = &accountDao{}

func (dao *accountDao) GetAbiByAddress(ctx context.Context, tx *ent.Client, address string) (a *ent.Account, err error) {
	return tx.Account.Get(ctx, address)
}

func (dao *accountDao) GetAccountByAddresses(ctx context.Context, tx *ent.Client, address []string) (accounts ent.Accounts, err error) {
	return tx.Account.Query().Where(
		account.IDIn(address...),
	).All(ctx)
}

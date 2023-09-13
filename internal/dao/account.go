package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/account"
	"github.com/BlockPILabs/aa-scan/internal/memo"
	"time"
)

type accountDao struct {
	baseDao
}

var AccountDao = &accountDao{}

func (dao *accountDao) GetAbiByAddress(ctx context.Context, tx *ent.Client, address string) (a *ent.Account, err error) {
	return tx.Account.Get(ctx, address)
}
func (dao *accountDao) GetAbiByAddressWithMemo(ctx context.Context, tx *ent.Client, address string) (a *ent.Account, err error) {
	key := "account:" + address
	memoAccount, b := memo.Get(key)
	if b {
		a, b = memoAccount.(*ent.Account)
		if b {
			return a, nil
		}
	}
	a, err = tx.Account.Get(ctx, address)
	if err != nil {
		return nil, err
	}
	memo.SetWithTTL(key, a, 10, time.Hour*30)
	return
}

func (dao *accountDao) GetAccountByAddresses(ctx context.Context, tx *ent.Client, address []string) (accounts ent.Accounts, err error) {
	return tx.Account.Query().Where(
		account.IDIn(address...),
	).All(ctx)
}

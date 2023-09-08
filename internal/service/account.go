package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/memo"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"time"
)

type accountService struct {
}

var AccountService = &accountService{}

func (*accountService) GetAbiByAddress(ctx context.Context, client *ent.Client, address string) (*abi.ABI, error) {
	ctx, logger := log.With(ctx, "service", "GetBlock")
	key := "addr:abi:" + address
	a, ok := memo.Get(key)
	if ok {
		ab, ok := a.(*abi.ABI)
		if ok && ab != nil {
			return ab, nil
		}
		return nil, errors.New("abi not found")
	}

	account, err := dao.AccountDao.GetAbiByAddress(ctx, client, address)
	if err != nil {
		memo.SetWithTTL(key, err, 1, time.Minute)
		logger.Error("GetAbiByAddress", "err", err)
		return nil, err
	}

	if !account.IsContract {
		memo.SetWithTTL(key, err, 1, time.Hour)
		return nil, errors.New("account is not construct")
	}

	if len(account.Abi) < 1 {
		memo.SetWithTTL(key, err, 1, time.Hour)
		return nil, errors.New("account abi not found")
	}

	aa, err := abi.JSON(bytes.NewBufferString(account.Abi))
	if err != nil {
		memo.SetWithTTL(key, err, 1, time.Minute*10)
		return nil, err
	}
	memo.SetWithTTL(key, aa, 1, time.Hour*24*30)
	return &aa, nil
}

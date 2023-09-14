package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type aaAccountService struct {
}

var AaAccountService = &aaAccountService{}

func (s *aaAccountService) GetAaAccountRecord(ctx context.Context, client *ent.Client, vo vo.AaAccountRequestVo) (*vo.AaAccountRecord, error) {
	ctx, logger := log.With(ctx, "service", "GetAaAccountRecord")
	logger.Info("GetAaAccountRecord ... ")

	address := vo.Address
	record, err := dao.AaAccountDao.GetAaAccountRecord(ctx, client, address)
	if err != nil {
		return nil, err
	}
	return record, nil
}

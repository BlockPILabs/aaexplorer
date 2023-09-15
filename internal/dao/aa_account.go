package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/userassetinfo"
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

	return query.Limit(50).All(ctx)
}

func (dao *aaAccountDao) GetAaAccountRecord(ctx context.Context, tx *ent.Client, address string) (*vo.AaAccountRecord, error) {
	record := vo.AaAccountRecord{}
	err := tx.AaAccountData.Query().Where(aaaccountdata.ID(address)).Modify(func(s *sql.Selector) {
		s.LeftJoin(sql.Table(userassetinfo.Table)).On(aaaccountdata.FieldID, userassetinfo.FieldAccountAddress)
	}).GroupBy(
		aaaccountdata.FieldID,
		aaaccountdata.FieldAaType,
		aaaccountdata.FieldFactory,
		aaaccountdata.FieldFactoryTime,
	).Aggregate(ent.As(ent.Sum(userassetinfo.FieldAmount), "total_amount")).Scan(ctx, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

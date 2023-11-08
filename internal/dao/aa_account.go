package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/userassetinfo"
	"github.com/BlockPILabs/aaexplorer/internal/utils"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type aaAccountDao struct {
	baseDao
}

var AaAccountDao = &aaAccountDao{}

func (*aaAccountDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		aaaccountdata.FieldUpdateTime,
		aaaccountdata.FieldFactoryTime,
		aaaccountdata.FieldUserOpsNum,
		aaaccountdata.FieldTotalBalanceUsd,
	}
}
func (dao *aaAccountDao) Sort(ctx context.Context, query *ent.AaAccountDataQuery, sort int, order int) *ent.AaAccountDataQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		f := dao.sortField(ctx, dao.GetSortFields(ctx), sort)
		switch f {
		case "", config.Default:
			query.Order(aaaccountdata.ByUpdateTime(opts...))
		default:
			query.Order(sql.OrderByField(f, opts...).ToFunc())
		}
	}
	return query
}
func (dao *aaAccountDao) Search(ctx context.Context, tx *ent.Client, req vo.SearchAllRequest) (a ent.AaAccountDataSlice, err error) {

	term := utils.Fix0x(strings.ToLower(req.Term))

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

type AaAccountScan struct {
	Address      *string
	Aa_type      *string
	Factory      *string
	Factory_time *time.Time
	Total_amount *decimal.Decimal
}

func (dao *aaAccountDao) GetAaAccountRecord(ctx context.Context, tx *ent.Client, address string) (*vo.AaAccountRecord, error) {
	var record []*AaAccountScan
	err := tx.AaAccountData.Query().Where(aaaccountdata.ID(address)).GroupBy(
		aaaccountdata.FieldID,
		aaaccountdata.FieldAaType,
		aaaccountdata.FieldFactory,
		aaaccountdata.FieldFactoryTime,
	).Aggregate(func(selector *sql.Selector) string {
		t := sql.Table(userassetinfo.Table)
		selector.LeftJoin(t).On(selector.C(aaaccountdata.FieldID), t.C(userassetinfo.FieldAccountAddress))
		return sql.As(sql.Sum(t.C(userassetinfo.FieldAmount)), "total_amount")
	}).Scan(ctx, &record)
	if err != nil {
		return nil, err
	}

	if len(record) <= 0 {
		return nil, nil
	}
	ret := &vo.AaAccountRecord{
		Address:     record[0].Address,
		AaType:      record[0].Aa_type,
		Factory:     record[0].Factory,
		TotalAmount: record[0].Total_amount,
	}
	if record[0].Factory_time != nil {
		ret.FactoryTime = record[0].Factory_time.UnixMilli()
	}
	return ret, nil
}

func (dao *aaAccountDao) Pagination(ctx context.Context, tx *ent.Client, req vo.GetAccountsRequest) (list ent.AaAccountDataSlice, total int, err error) {
	query := tx.AaAccountData.Query()

	if len(req.Address) > 0 {
		query = query.Where(
			aaaccountdata.IDEQ(req.Address),
		)
	}
	if len(req.Factory) > 0 {
		query = query.Where(
			aaaccountdata.FactoryEQ(req.Factory),
		)
	}

	if req.TotalCount > 0 {
		total = req.TotalCount
	} else {
		total = query.CountX(ctx)
	}

	if total < 1 || req.GetOffset() > total {
		return
	}

	// sort
	query = dao.Sort(ctx, query, req.Sort, req.Order)
	// limit
	query = query.WithAccount().
		Offset(req.GetOffset()).
		Limit(req.PerPage)

	list, err = query.All(ctx)
	return
}

func (dao *aaAccountDao) AaAccountExists(ctx context.Context, tx *ent.Client, address string) bool {
	exist, err := tx.AaAccountData.Query().Exist(ctx)
	if err != nil {
		return false
	}
	return exist
}

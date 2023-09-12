package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"golang.org/x/sync/errgroup"
	"strconv"
)

type searchService struct {
}

var SearchService = &searchService{}

func (s *searchService) SearchAll(ctx context.Context, req vo.SearchAllRequest) (*vo.SearchAllResponse, error) {
	res := &vo.SearchAllResponse{
		WalletAccounts: []*vo.SearchAllAccount{},
		Paymasters:     []*vo.SearchAllAccount{},
		Bundlers:       []*vo.SearchAllAccount{},
		Transactions:   []*vo.SearchAllTransaction{},
		Blocks:         []*vo.SearchAllBlock{},
	}

	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return res, err
	}
	term := req.Term
	if utils.Has0xPrefix(term) {
		term = term[2:]
	}
	maxResult := 10
	wg, ctx := errgroup.WithContext(ctx)

	if utils.IsHexSting(term) {
		wg.Go(func() error {
			accounts, err := dao.AaAccountDao.Search(ctx, client, req)
			if err != nil {
				log.Context(ctx).Warn("search account error", "err", err, "term", req.Term)
				return nil
			}
			for _, account := range accounts {
				sa := &vo.SearchAllAccount{Address: account.ID, AaType: account.AaType}
				switch account.AaType {
				case config.AaAccountTypePaymaster:
					if len(res.Paymasters) < maxResult {
						res.Paymasters = append(res.Paymasters, sa)
					}
				case config.AaAccountTypeBundler:
					if len(res.Paymasters) < maxResult {
						res.Bundlers = append(res.Bundlers, sa)
					}
				default:
					if len(res.Paymasters) < maxResult {
						res.WalletAccounts = append(res.WalletAccounts, sa)
					}
				}
			}
			return nil
		})

		wg.Go(func() error {
			blocks, err := dao.BlockDao.Search(ctx, client, req.Term)
			if err != nil {
				log.Context(ctx).Warn("search block error", "err", err, "req.Term", req.Term)
				return nil
			}
			for _, block := range blocks {
				res.Blocks = append(res.Blocks, &vo.SearchAllBlock{BlockNumber: block.ID})
			}
			return nil
		})
	}

	if utils.IsNumber(req.Term) {
		parseInt, _ := strconv.ParseInt(req.Term, 10, 64)
		if parseInt > 0 {
			wg.Go(func() error {
				block, err := dao.BlockDao.GetByBlockNumber(ctx, client, parseInt)
				if err != nil {
					log.Context(ctx).Warn("get block by number error", "err", err, "req.Term", req.Term, "number", parseInt)
					return nil
				}
				res.Blocks = append(res.Blocks, &vo.SearchAllBlock{BlockNumber: block.ID})
				return nil
			})
		}
	}

	err = wg.Wait()
	if err != nil {
		return nil, err
	}
	return res, nil
}

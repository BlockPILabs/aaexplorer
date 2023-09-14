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
		//WalletAccounts: []*vo.SearchAllAccount{},
		//Paymasters:     []*vo.SearchAllAccount{},
		//Bundlers:       []*vo.SearchAllAccount{},
		//Transactions:   []*vo.SearchAllTransaction{},
		//Blocks:         []*vo.SearchAllBlock{},
		Data: make([]*vo.SearchAllResponseData, 0),
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

	var paymasters []*vo.SearchAllAccount
	var bundlers []*vo.SearchAllAccount
	var walletAccounts []*vo.SearchAllAccount
	var blocks []*vo.SearchAllBlock
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
					if len(paymasters) < maxResult {
						paymasters = append(paymasters, sa)
					}
				case config.AaAccountTypeBundler:
					if len(bundlers) < maxResult {
						bundlers = append(bundlers, sa)
					}
				default:
					if len(walletAccounts) < maxResult {
						walletAccounts = append(walletAccounts, sa)
					}
				}
			}
			return nil
		})

		wg.Go(func() error {
			_blocks, err := dao.BlockDao.Search(ctx, client, req.Term)
			if err != nil {
				log.Context(ctx).Warn("search block error", "err", err, "req.Term", req.Term)
				return nil
			}
			for _, block := range _blocks {
				blocks = append(blocks, &vo.SearchAllBlock{BlockNumber: block.ID, BlockHash: block.Hash})
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
				blocks = append(blocks, &vo.SearchAllBlock{BlockNumber: block.ID, BlockHash: block.Hash})
				return nil
			})
		}
	}

	err = wg.Wait()

	if len(walletAccounts) > 0 {
		res.Data = append(res.Data,
			&vo.SearchAllResponseData{
				Type:    "Wallet account",
				Records: walletAccounts,
			},
		)
	}

	if len(paymasters) > 0 {
		res.Data = append(res.Data,
			&vo.SearchAllResponseData{
				Type:    "Paymaster",
				Records: paymasters,
			},
		)
	}

	if len(bundlers) > 0 {
		res.Data = append(res.Data,
			&vo.SearchAllResponseData{
				Type:    "Bundler",
				Records: bundlers,
			},
		)
	}

	if len(blocks) > 0 {
		res.Data = append(res.Data,
			&vo.SearchAllResponseData{
				Type:    "Block",
				Records: blocks,
			},
		)
	}

	if err != nil {
		return nil, err
	}
	return res, nil
}

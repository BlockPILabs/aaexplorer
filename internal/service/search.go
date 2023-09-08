package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/internal/vo"
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
	if utils.IsHexSting(term) {
		accounts, _ := dao.AaAccountDao.Search(ctx, client, req)
		for _, account := range accounts {
			sa := &vo.SearchAllAccount{Address: account.ID}
			switch account.AaType {
			case config.AaAccountTypePaymaster:
				res.Paymasters = append(res.Paymasters, sa)
			case config.AaAccountTypeBundler:
				res.Bundlers = append(res.Bundlers, sa)
			default:
				res.WalletAccounts = append(res.WalletAccounts, sa)
			}
		}
	}

	return res, nil
}

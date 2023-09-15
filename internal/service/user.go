package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/BlockPILabs/aa-scan/service"
	"github.com/shopspring/decimal"
)

func GetUserBalance(ctx context.Context, req vo.UserBalanceRequest) (*vo.UserBalanceResponse, error) {
	network := req.Network
	var resp = &vo.UserBalanceResponse{}

	account := req.AccountAddress
	balanceDetails := service.GetWalletBalanceDetail(account, network)
	if balanceDetails == nil {
		resp.TotalUsd = decimal.Zero
		return resp, nil
	}
	var totalUsd = decimal.Zero
	var details []*vo.BalanceDetail
	for _, detail := range balanceDetails {
		balance := &vo.BalanceDetail{
			TokenAddress:  detail.ContractAddress,
			TokenAmount:   detail.Amount,
			Percentage:    detail.Percent,
			TokenValueUsd: detail.ValueUsd.RoundDown(4),
		}
		details = append(details, balance)
		totalUsd = totalUsd.Add(detail.ValueUsd).RoundDown(4)
	}
	resp.BalanceDetails = details
	resp.TotalUsd = totalUsd

	return resp, nil
}

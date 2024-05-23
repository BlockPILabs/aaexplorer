package dao

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaasset"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/monitor"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/shopspring/decimal"
	"strings"
)

type monitorDao struct {
	baseDao
}

var MonitorDao = &monitorDao{}

func (dao *monitorDao) ListMonitorDao(ctx context.Context, req vo.ListWatchingAddressRequest) ([]*vo.WatchingAddress, int, error) {
	client, err := entity.Client(ctx)

	if err != nil {
		return nil, 0, err
	}

	query := client.Monitor.Query().Where(monitor.UserAddressEQ(req.UserAddress))

	total := query.CountX(ctx)

	monitorList, err := query.Offset(req.GetOffset()).Limit(req.GetPerPage()).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	if len(monitorList) == 0 {
		return nil, 0, nil
	}

	client, err = entity.Client(ctx, req.Network)
	if err != nil {
		return nil, 0, err
	}

	var list []*vo.WatchingAddress

	for _, monitorItem := range monitorList {
		monitorAsset, err := client.AaAsset.Query().Where(aaasset.IDEQ(monitorItem.MonitorAddress)).All(ctx)
		if err != nil {
			continue
		}
		monitorBalance := decimal.Zero
		if len(monitorAsset) > 0 {
			monitorBalance = monitorAsset[0].AssetValue
		}
		monitorAddress := monitorItem.MonitorAddress
		switch monitorItem.MonitorAddressType {
		case "bundler":
			bundlerInfo, err := client.BundlerInfo.Query().Where(bundlerinfo.IDEqualFold(monitorAddress)).All(ctx)
			if err != nil {
				continue
			}
			profits := decimal.Zero
			userOpsNum := int64(0)
			if len(bundlerInfo) > 0 {
				profits = bundlerInfo[0].FeeEarnedD1
				userOpsNum = bundlerInfo[0].UserOpsNum
			}

			list = append(list, &vo.WatchingAddress{
				Network:         req.Network,
				AddressType:     "Bundler",
				MonitorAddress:  monitorItem.MonitorAddress,
				Balance:         monitorBalance,
				Profits24H:      profits,
				SponsoredGas24H: decimal.Zero,
				TotalUserOps:    userOpsNum,
			})
		case "paymaster":
			paymaster, err := client.PaymasterInfo.Query().Where(paymasterinfo.IDEQ(monitorItem.MonitorAddress)).All(ctx)
			if err != nil {
				continue
			}
			gasSponsored := decimal.Zero
			userOpsNum := int64(0)
			if len(paymaster) > 0 {
				gasSponsored = paymaster[0].GasSponsoredD1
				userOpsNum = paymaster[0].UserOpsNum
			}

			list = append(list, &vo.WatchingAddress{
				Network:         req.Network,
				AddressType:     "Paymaster",
				MonitorAddress:  monitorItem.MonitorAddress,
				Balance:         monitorBalance,
				Profits24H:      decimal.Zero,
				SponsoredGas24H: gasSponsored,
				TotalUserOps:    userOpsNum,
			})
		case "aa":
			list = append(list, &vo.WatchingAddress{
				Network:         req.Network,
				AddressType:     "Account",
				MonitorAddress:  monitorItem.MonitorAddress,
				Balance:         monitorBalance,
				Profits24H:      decimal.Zero,
				SponsoredGas24H: decimal.Zero,
				TotalUserOps:    int64(0),
			})
		case "entry_point":
			list = append(list, &vo.WatchingAddress{
				Network:         req.Network,
				AddressType:     "Contract Account",
				MonitorAddress:  monitorItem.MonitorAddress,
				Balance:         monitorBalance,
				Profits24H:      decimal.Zero,
				SponsoredGas24H: decimal.Zero,
				TotalUserOps:    int64(0),
			})
		default:
			list = append(list, &vo.WatchingAddress{
				Network:         req.Network,
				AddressType:     strings.ToUpper(monitorItem.MonitorAddressType[:1]) + monitorItem.MonitorAddressType[1:],
				MonitorAddress:  monitorItem.MonitorAddress,
				Balance:         monitorBalance,
				Profits24H:      decimal.Zero,
				SponsoredGas24H: decimal.Zero,
				TotalUserOps:    int64(0),
			})
		}
	}

	return list, total, nil
}

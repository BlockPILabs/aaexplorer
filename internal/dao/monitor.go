package dao

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaasset"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/monitor"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
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
		monitorAsset, err := client.AaAsset.Query().Where(aaasset.IDEQ(monitorItem.MonitorAddress)).Only(ctx)
		if err != nil {
			continue
		}

		monitorBalance, _ := monitorAsset.AssetValue.Float64()
		switch monitorItem.MonitorAddressType {
		case "bundler":
			bundlerInfo, err := client.BundlerInfo.Query().Where(bundlerinfo.IDEQ(monitorItem.MonitorAddress)).Only(ctx)
			if err != nil {
				continue
			}
			profits, _ := bundlerInfo.FeeEarnedD1.Float64()

			list = append(list, &vo.WatchingAddress{
				Network:         req.Network,
				AddressType:     "Bundler",
				MonitorAddress:  monitorItem.MonitorAddress,
				Balance:         monitorBalance,
				Profits24H:      profits,
				SponsoredGas24H: float64(0),
				TotalUserOps:    bundlerInfo.UserOpsNum,
			})
		case "paymaster":
			paymaster, err := client.PaymasterInfo.Query().Where(paymasterinfo.IDEQ(monitorItem.MonitorAddress)).Only(ctx)
			if err != nil {
				continue
			}

			gasSponsored, _ := paymaster.GasSponsoredD1.Float64()

			list = append(list, &vo.WatchingAddress{
				Network:         req.Network,
				AddressType:     "Paymaster",
				MonitorAddress:  monitorItem.MonitorAddress,
				Balance:         monitorBalance,
				Profits24H:      float64(0),
				SponsoredGas24H: gasSponsored,
				TotalUserOps:    paymaster.UserOpsNum,
			})
		default:
			list = append(list, &vo.WatchingAddress{
				Network:         req.Network,
				AddressType:     strings.ToTitle(monitorItem.MonitorAddressType),
				MonitorAddress:  monitorItem.MonitorAddress,
				Balance:         monitorBalance,
				Profits24H:      float64(0),
				SponsoredGas24H: float64(0),
				TotalUserOps:    int64(0),
			})
		}
	}

	return list, total, nil
}

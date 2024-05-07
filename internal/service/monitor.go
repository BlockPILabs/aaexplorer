package service

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/monitor"
	interlog "github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/BlockPILabs/aaexplorer/util"
	"strings"
	"time"
)

var logger = interlog.L()

func SetLogger(lg interlog.Logger) {
	logger = lg
}

func AddMonitor(ctx context.Context, req vo.AddMonitorRequest) (*vo.AddMonitorResponse, error) {

	monitorAddress := req.MonitorAddress
	if len(monitorAddress) == 0 {
		return nil, nil
	}
	network := req.Network
	monitorAddress = strings.ToLower(monitorAddress)
	client, err := entity.Client(ctx)
	if err != nil {
		return nil, err
	}
	var resp = &vo.AddMonitorResponse{}

	userAddress, err := util.RecoverSignFromMetamask(req.Message, req.Sign)
	if err != nil {
		logger.Error("AddMonitor sign err ", "monitorAddress", monitorAddress, "network", network, "msg", err)
		return nil, err
	}
	userAddress = strings.ToLower(userAddress)

	olds, err := client.Monitor.Query().Where(monitor.MonitorAddressEqualFold(monitorAddress), monitor.UserAddressEqualFold(userAddress)).All(ctx)
	if len(olds) > 0 {
		return nil, nil
	}

	nClient, err := entity.Client(ctx, network)
	if err != nil {
		return nil, nil
	}
	aaAccounts, err := nClient.AaAccountData.Query().Where(aaaccountdata.IDEqualFold(monitorAddress)).All(ctx)
	var monitorType = ""
	if len(aaAccounts) > 0 {
		monitorType = aaAccounts[0].AaType
	}

	monitor := client.Monitor.Create().SetCreateTime(time.Now()).SetMonitorAddress(monitorAddress).SetUserAddress(userAddress).SetMonitorAddressType(monitorType)
	_, err = monitor.Save(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func RemoveMonitor(ctx context.Context, req vo.RemoveMonitorRequest) (*vo.RemoveMonitorResponse, error) {

	monitorAddress := req.MonitorAddress
	if len(monitorAddress) == 0 {
		return nil, nil
	}
	network := req.Network
	monitorAddress = strings.ToLower(monitorAddress)
	client, err := entity.Client(ctx)
	if err != nil {
		return nil, err
	}
	var resp = &vo.RemoveMonitorResponse{}

	userAddress, err := util.RecoverSignFromMetamask(req.Message, req.Sign)
	if err != nil {
		logger.Error("RemoveMonitor sign err ", "monitorAddress", monitorAddress, "network", network, "msg", err)
		return nil, err
	}
	userAddress = strings.ToLower(userAddress)

	olds, err := client.Monitor.Query().Where(monitor.MonitorAddressEqualFold(monitorAddress), monitor.UserAddressEqualFold(userAddress)).All(ctx)
	if len(olds) == 0 {
		logger.Info("RemoveMonitor not exist ", "userAddress", userAddress, "monitorAddress", monitorAddress)
		return nil, nil
	}

	client.Monitor.Delete().Where(monitor.UserAddressEqualFold(userAddress), monitor.UserAddressEqualFold(monitorAddress)).Exec(ctx)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

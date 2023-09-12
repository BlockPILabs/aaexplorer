package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aacontractinteract"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/useroptypestatistic"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/shopspring/decimal"
	"log"
)

func GetUserOpType(ctx context.Context, req vo.UserOpsTypeRequest) (*vo.UserOpsTypeResponse, error) {
	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	timeRange := req.TimeRange
	var resp = &vo.UserOpsTypeResponse{}

	userOpTypes, err := client.UserOpTypeStatistic.Query().Where(useroptypestatistic.StatisticTypeEQ(timeRange), useroptypestatistic.NetworkEqualFold(network)).All(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp = getUserOpTypeResponse(userOpTypes)

	return resp, nil
}

func getUserOpTypeResponse(types []*ent.UserOpTypeStatistic) *vo.UserOpsTypeResponse {
	if len(types) == 0 {
		return nil
	}

	var totalNum = int64(0)
	for _, opType := range types {
		totalNum += opType.OpNum
	}

	var resp = &vo.UserOpsTypeResponse{}
	var userOpsInfos []*vo.UserOpsType
	for _, opType := range types {
		userOpType := &vo.UserOpsType{
			UserOpType: opType.UserOpType,
			Rate:       decimal.NewFromInt(opType.OpNum).DivRound(decimal.NewFromInt(totalNum), 4),
		}
		userOpsInfos = append(userOpsInfos, userOpType)
	}
	resp.UserOpTypes = userOpsInfos
	return resp
}

func GetAAContractInteract(ctx context.Context, req vo.AAContractInteractRequest) (*vo.AAContractInteractResponse, error) {
	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	timeRange := req.TimeRange
	var resp = &vo.AAContractInteractResponse{}

	contractInteract, err := client.AAContractInteract.Query().Where(aacontractinteract.StatisticTypeEQ(timeRange), aacontractinteract.NetworkEqualFold(network)).All(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp = getAAContractInteractResponse(contractInteract)

	return resp, nil
}

func getAAContractInteractResponse(interacts []*ent.AAContractInteract) *vo.AAContractInteractResponse {
	if len(interacts) == 0 {
		return nil
	}

	var totalNum = int64(0)
	for _, interact := range interacts {
		totalNum += interact.InteractNum
	}

	var resp = &vo.AAContractInteractResponse{}
	var contractInteractArr []*vo.AAContractInteract
	for _, interact := range interacts {
		contractInteract := &vo.AAContractInteract{
			ContractAddress: interact.ContractAddress,
			Rate:            decimal.NewFromInt(interact.InteractNum).DivRound(decimal.NewFromInt(totalNum), 4),
			SingleNum:       interact.InteractNum,
		}
		contractInteractArr = append(contractInteractArr, contractInteract)
	}
	resp.AAContractInteract = contractInteractArr
	resp.TotalNum = totalNum
	return resp
}

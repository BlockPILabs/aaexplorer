package service

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aacontractinteract"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/functionsignature"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transfertransaction"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/useroptypestatistic"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/shopspring/decimal"
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
		if len(opType.UserOpType) == 0 {
			continue
		}
		totalNum += opType.OpNum
	}

	var resp = &vo.UserOpsTypeResponse{}
	var userOpsInfos []*vo.UserOpsType
	for _, opType := range types {
		if len(opType.UserOpType) == 0 {
			continue
		}
		userOpType := &vo.UserOpsType{
			UserOpType: opType.UserOpType,
			Rate:       decimal.NewFromInt(opType.OpNum).DivRound(decimal.NewFromInt(totalNum), 4),
		}
		userOpsInfos = append(userOpsInfos, userOpType)
	}
	sort.Sort(vo.ByUserOpsTypeNum(userOpsInfos))
	var finalOpsInfos []*vo.UserOpsType
	var totalRate = decimal.Zero
	client, err := entity.Client(context.Background())
	if err != nil {
		return nil
	}

	for i, userOpsInfo := range userOpsInfos {
		if i >= config.AnalyzeTop7 {
			break
		}
		funcSig, _ := client.FunctionSignature.Query().Where(functionsignature.IDEqualFold("0x" + userOpsInfo.UserOpType)).All(context.Background())
		if len(funcSig) > 0 {
			userOpsInfo.UserOpType = funcSig[0].Name
		}
		finalOpsInfos = append(finalOpsInfos, userOpsInfo)
		totalRate = totalRate.Add(userOpsInfo.Rate)
	}
	if totalRate.Cmp(decimal.NewFromInt(1)) > 0 {
		last := finalOpsInfos[len(finalOpsInfos)-1]
		last.UserOpType = "other"
		last.Rate = decimal.NewFromInt(1).Sub(totalRate.Sub(last.Rate))
	} else {
		leftRate := decimal.NewFromInt(1).Sub(totalRate)
		if leftRate.Cmp(decimal.Zero) > 0 {
			newLast := &vo.UserOpsType{
				UserOpType: "other",
				Rate:       decimal.NewFromInt(1).Sub(totalRate),
			}
			finalOpsInfos = append(finalOpsInfos, newLast)
		}

	}

	resp.UserOpTypes = finalOpsInfos
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
		if len(interact.ContractAddress) == 0 {
			continue
		}
		totalNum += interact.InteractNum
	}

	var resp = &vo.AAContractInteractResponse{}
	var contractInteractArr []*vo.AAContractInteract
	for _, interact := range interacts {
		if len(interact.ContractAddress) == 0 {
			continue
		}
		contractInteract := &vo.AAContractInteract{
			ContractAddress: interact.ContractAddress,
			Rate:            decimal.NewFromInt(interact.InteractNum).DivRound(decimal.NewFromInt(totalNum), 4),
			SingleNum:       interact.InteractNum,
		}
		contractInteractArr = append(contractInteractArr, contractInteract)
	}
	sort.Sort(vo.ByContractNum(contractInteractArr))
	var finalContractInteracts []*vo.AAContractInteract
	var totalRate = decimal.Zero
	var topNum = int64(0)
	for i, contractInteract := range contractInteractArr {
		if i >= config.AnalyzeTop7 {
			break
		}
		finalContractInteracts = append(finalContractInteracts, contractInteract)
		totalRate = totalRate.Add(contractInteract.Rate)
		topNum += contractInteract.SingleNum
	}
	if totalRate.Cmp(decimal.NewFromInt(1)) > 0 {
		last := finalContractInteracts[len(finalContractInteracts)-1]
		last.ContractAddress = "other"
		last.Rate = decimal.NewFromInt(1).Sub(totalRate.Sub(last.Rate))
	} else {
		leftRate := decimal.NewFromInt(1).Sub(totalRate)
		if leftRate.Cmp(decimal.Zero) > 0 {
			newLast := &vo.AAContractInteract{
				ContractAddress: "other",
				Rate:            decimal.NewFromInt(1).Sub(totalRate),
				SingleNum:       totalNum - topNum,
			}
			finalContractInteracts = append(finalContractInteracts, newLast)
		}

	}
	resp.AAContractInteract = finalContractInteracts
	resp.TotalNum = totalNum
	return resp
}

func GetHotAAToken(ctx context.Context, req vo.HotAARequest) (*vo.HotAAResponse, error) {
	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	var resp = &vo.HotAAResponse{}

	day1 := time.UnixMilli(time.Now().UnixMilli() - 24*3600*1000)
	//Where(transfertransaction.TimeGTE(day1))
	var results []*HotAA
	err = client.TransferTransaction.Query().Where(transfertransaction.TimeGTE(day1)).GroupBy(transfertransaction.FieldTokenSymbol).Aggregate(ent.Count()).Scan(ctx, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}

	sort.Sort(ByHotAA(results))
	var details []vo.TokenDetail
	for idx, res := range results {
		if idx > 14 {
			break
		}
		detail := vo.TokenDetail{
			TokenSymbol: res.TokenSymbol,
			Count:       res.Count,
		}
		details = append(details, detail)

	}

	resp.TokenDetails = details

	return resp, nil
}

type HotAA struct {
	TokenSymbol string `json:"token_symbol"`
	Count       int64  `json:"count"`
}

type ByHotAA []*HotAA

func (b ByHotAA) Len() int      { return len(b) }
func (b ByHotAA) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByHotAA) Less(i, j int) bool {
	return b[i].Count-b[j].Count > 0
}

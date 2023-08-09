package moralis

import (
	"context"
	"encoding/json"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/predicate"
	"github.com/shopspring/decimal"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const ApiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6IjIxZTBkZmU5LTlkNTItNGQ5ZC05YmQzLTdhZjBhYjAyNDFhYiIsIm9yZ0lkIjoiMzUwOTE1IiwidXNlcklkIjoiMzYwNjcwIiwidHlwZUlkIjoiY2VhNjZmN2MtNTYwMi00NGQzLWE5YzUtNDhjMjA5MmQzNzU5IiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE2OTA3ODgyMDEsImV4cCI6NDg0NjU0ODIwMX0.VxVXZ6z9y3QY9_JfEsQvxBcs2SmFk05OObAptKofGxc"
const MoralisUrl = "https://deep-index.moralis.io/api/v2"

type TokenBalance struct {
	TokenAddress string          `json:"token_address"`
	Symbol       string          `json:"symbol"`
	Name         string          `json:"name"`
	Logo         string          `json:"logo"`
	Thumbnail    string          `json:"thumbnail"`
	Decimals     int32           `json:"decimals"`
	Balance      decimal.Decimal `json:"balance,string"`
	PossibleSpam bool            `json:"possible_spam"`
}

type TokenPrice struct {
	TokenName        string          `json:"tokenName"`
	TokenSymbol      string          `json:"tokenSymbol"`
	TokenLogo        string          `json:"tokenLogo"`
	TokenDecimals    string          `json:"tokenDecimals"`
	NativePrice      NativePrice     `json:"nativePrice"`
	UsdPrice         decimal.Decimal `json:"usdPrice"`
	PercentChange24h string          `json:"24hrPercentChange"`
	ExchangeAddress  string          `json:"exchangeAddress"`
	ExchangeName     string          `json:"exchangeName"`
	TokenAddress     string          `json:"tokenAddress"`
}

type NativePrice struct {
	Value    string
	Decimals int32
	Name     string
	Symbol   string
}

type tokensLoad struct {
	TokenAddress string `json:"token_address"`
}

type pricePayload struct {
	Tokens []*tokensLoad `json:"tokens"`
}

func GetTokenBalance(address string, network string) []*TokenBalance {

	url := MoralisUrl + "/" + address + "/erc20?chain=" + network

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-Key", ApiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	data := string(body)

	//fmt.Println(res)
	//fmt.Println(data)
	var tokenBalanceArr []*TokenBalance

	err := json.Unmarshal([]byte(data), &tokenBalanceArr)
	if err != nil {
		fmt.Println("parse json:", err)
		return nil
	}

	var tokens []string
	for _, t := range tokenBalanceArr {
		fmt.Println(json.Marshal(t))
		tokens = append(tokens, t.TokenAddress)
	}

	GetTokenPriceBatch(tokens)
	fmt.Println("aaaaa")

	return tokenBalanceArr

}

func GetTokenPriceBatch(tokens []string) []*TokenPrice {

	url := MoralisUrl + "/erc20/prices?chain=eth&include=percent_change"
	var tokenArr []*tokensLoad
	for _, tokenAddr := range tokens {
		tokenArr = append(tokenArr, &tokensLoad{TokenAddress: tokenAddr})
	}
	priceLoad := pricePayload{Tokens: tokenArr}
	jsonData, err := json.Marshal(priceLoad)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	payload := strings.NewReader(string(jsonData))

	req, _ := http.NewRequest("GET", url, payload)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", ApiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	data := string(body)
	fmt.Println(data)

	var tokenPriceArr []*TokenPrice

	err = json.Unmarshal([]byte(data), &tokenPriceArr)
	if err != nil {
		fmt.Println("parse json:", err)
		return nil
	}

	return tokenPriceArr

}

func GetTokenPrice(token string, network string) *TokenPrice {

	url := MoralisUrl + "/erc20/" + token + "/price?chain=" + network + "&include=percent_change"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", ApiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	data := string(body)
	fmt.Println(data)

	var tokenPrice *TokenPrice

	err := json.Unmarshal([]byte(data), &tokenPrice)
	if err != nil {
		fmt.Println("parse json:", err)
		return nil
	}

	return tokenPrice

}

type TokenPriceInfo struct {
	Id              int64
	ContractAddress string
	Symbol          string
	TokenPrice      decimal.Decimal
	LastTime        int64
	CreateTime      time.Time
	UpdateTime      time.Time
}

func syncPrice(contractAddress string, network string) {
	var tokenInfo TokenPriceInfo
	lastTime := tokenInfo.LastTime
	curTime := time.Now().UnixNano() / 1e6
	if curTime-lastTime < 5*60*1000 {
		return
	}
	price := GetTokenPrice(contractAddress, network)
	updateToken := &TokenPriceInfo{
		Id:         tokenInfo.Id,
		TokenPrice: price.UsdPrice,
		LastTime:   curTime,
		UpdateTime: time.Now(),
	}
	time1 := updateToken.UpdateTime
	fmt.Println(time1)
}

type UserAssetInfo struct {
	Id              int64
	AccountAddress  string
	ContractAddress string
	Symbol          string
	Network         string
	Amount          decimal.Decimal
	LastTime        int64
	CreateTime      time.Time
}

type WalletBalanceResp struct {
	ContractAddress string
	Symbol          string
	Value           string
	Percent         string
}

func GetUserAsset(accountAddress string, network string) []*UserAssetInfo {
	var userAssetInfos []*UserAssetInfo
	if len(userAssetInfos) == 0 {
		return nil
	}
	lastTime := userAssetInfos[0].LastTime
	curTime := time.Now().UnixNano() / 1e6
	if curTime-lastTime < 5*60*1000 {
		return userAssetInfos
	}
	tokenBalances := GetTokenBalance(accountAddress, network)
	//TODO change last time
	if len(tokenBalances) == 0 {
		return userAssetInfos
	}
	var userAssetInfoCreates []*ent.UserAssetInfoCreate
	client, err := entity.Client(context.Background())
	if err != nil {
		return nil
	}
	for _, tokenBalance := range tokenBalances {
		userAssetCreate := client.UserAssetInfo.Create().
			SetAccountAddress(accountAddress).
			SetContractAddress(tokenBalance.TokenAddress).
			SetSymbol(tokenBalance.Symbol).
			SetNetwork(network).
			SetAmount(tokenBalance.Balance).
			SetLastTime(curTime)
		userAssetInfoCreates = append(userAssetInfoCreates, userAssetCreate)
	}

	client.UserAssetInfo.Delete().Where(predicate.UserAssetInfo(sql.FieldEQ("account_address", accountAddress)))

	_, err = client.UserAssetInfo.CreateBulk(userAssetInfoCreates...).Save(context.Background())
	if err != nil {
		log.Println(err)
	}
	return userAssetInfos

}

func getMEV() {
	//1.查询失败的交易
	//2.查询相同sender，相同nonce的交易，查询出成功的
	//3.如果有，先看from，from不一样，再对比区块，区块差小于等于1，则是被抢跑的，标记处MEV
	//4.

}

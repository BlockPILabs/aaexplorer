package moralis

import (
	"context"
	"encoding/json"
	"entgo.io/ent/dialect/sql"
	"fmt"
	cfg "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/predicate"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/userassetinfo"
	"github.com/shopspring/decimal"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

const MoralisUrl = "https://deep-index.moralis.io/api/v2.2"

var config = cfg.DefaultConfig()

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

type NativeTokenBalance struct {
	Balance decimal.Decimal
}

func SetConfig(conf *cfg.Config) {
	config = conf
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
	network = strings.ToLower(network)
	if network == "ethereum" {
		network = "eth"
	}
	url := MoralisUrl + "/" + address + "/erc20?chain=" + network

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-Key", config.MoralisApiKey)

	res, _ := http.DefaultClient.Do(req)
	checkStatus(res)
	if res == nil {
		log.Printf("GetTokenBalance err, address: %s, network: %s", address, network)

	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	data := string(body)

	var tokenBalanceArr []*TokenBalance

	err := json.Unmarshal([]byte(data), &tokenBalanceArr)
	if err != nil {
		fmt.Println("parse json:", err)
		return nil
	}

	var tokens []string
	for _, t := range tokenBalanceArr {
		tokens = append(tokens, t.TokenAddress)
	}

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
		log.Println(err)
		return nil
	}
	payload := strings.NewReader(string(jsonData))

	req, _ := http.NewRequest("GET", url, payload)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", config.MoralisApiKey)

	res, _ := http.DefaultClient.Do(req)
	checkStatus(res)
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
	network = strings.ToLower(network)
	if network == "ethereum" {
		network = "eth"
	}
	url := MoralisUrl + "/erc20/" + token + "/price?chain=" + network + "&include=percent_change"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", config.MoralisApiKey)

	res, _ := http.DefaultClient.Do(req)
	checkStatus(res)
	if res == nil {
		log.Printf("GetTokenPrice err, token: %s, network: %s", token, network)
		return nil
	}

	if res.StatusCode != 200 {
		return nil
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	data := string(body)

	var tokenPrice *TokenPrice

	err := json.Unmarshal([]byte(data), &tokenPrice)
	if err != nil {
		fmt.Println("parse json:", err)
		return nil
	}

	return tokenPrice

}

func GetNativeTokenBalance(accountAddress string, network string) decimal.Decimal {
	if len(accountAddress) == 0 || len(network) == 0 {
		return decimal.Zero
	}
	network = strings.ToLower(network)
	if network == "ethereum" {
		network = "eth"
	}
	url := MoralisUrl + "/" + accountAddress + "/balance?chain=" + network

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-Key", config.MoralisApiKey)

	res, _ := http.DefaultClient.Do(req)
	checkStatus(res)
	if res == nil {
		log.Printf("GetNativeTokenBalance err, accountAddress: %s, network: %s", accountAddress, network)
		return decimal.Zero
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	data := string(body)

	var nativeTokenBalance *NativeTokenBalance

	err := json.Unmarshal([]byte(data), &nativeTokenBalance)
	if err != nil {
		return decimal.Zero
	}

	return nativeTokenBalance.Balance.DivRound(decimal.New(int64(math.Pow10(18)), 0), 18)
}

func checkStatus(res *http.Response) {
	if res == nil {
		log.Printf("call morialis err")
		return
	}
	code := res.StatusCode
	status := res.Status
	if code != 200 {
		log.Printf("call morialis failed, %s, %s", status, res.Request.URL)
	}
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

func GetUserAsset(accountAddress string, network string) []*ent.UserAssetInfo {
	client, err := entity.Client(context.Background())
	if err != nil {
		return nil
	}
	userAssetInfos, err := client.UserAssetInfo.Query().Where(userassetinfo.AccountAddressEqualFold(accountAddress), userassetinfo.NetworkEqualFold(network)).All(context.Background())
	if err != nil {
		return nil
	}
	curTime := time.Now().UnixMilli()
	if len(userAssetInfos) != 0 {
		lastTime := userAssetInfos[0].LastTime
		if curTime-lastTime < cfg.AssetExpireTime {
			return userAssetInfos
		}
	}

	native := GetNativeTokenBalance(accountAddress, network)
	tokenBalances := GetTokenBalance(accountAddress, network)
	if len(tokenBalances) == 0 && native == decimal.Zero {
		return userAssetInfos
	}
	if native != decimal.Zero {
		nativeToken := &TokenBalance{
			TokenAddress: cfg.ZeroAddress,
			Symbol:       GetNativeName(network),
			Name:         GetNativeName(network),
			Logo:         "",
			Thumbnail:    "",
			Decimals:     cfg.EvmDecimal,
			Balance:      native.DivRound(decimal.NewFromFloat(math.Pow10(18)), 18),
		}
		if tokenBalances == nil {
			tokenBalances = []*TokenBalance{}
		}
		tokenBalances = append(tokenBalances, nativeToken)
	}
	var userAssetInfoCreates []*ent.UserAssetInfoCreate
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

func GetNativeName(network string) string {

	if network == cfg.EthNetwork {
		return cfg.EthNative
	} else if network == cfg.BscNetwork {
		return cfg.BscNative
	} else if network == cfg.PolygonNetwork {
		return cfg.PolygonNative
	}

	return ""
}

package cmc

import (
	"encoding/json"
	interlog "github.com/BlockPILabs/aaexplorer/internal/log"
	"io"
	"net/http"
	"strconv"
)

const CmcUrl = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/map?sort=cmc_rank"
const CmcAuthVal = "e7f31052-0284-4ad6-8115-e35f2ecc8028"

var logger = interlog.L()

func SetLogger(lg interlog.Logger) {
	logger = lg
}

type TokenResp struct {
	Status Status      `json:"status"`
	Data   []TokenInfo `json:"data"`
}

type Status struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int64  `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int64  `json:"elapsed"`
	CreditCount  int64  `json:"credit_count"`
	Notice       string `json:"notice"`
}

type TokenInfo struct {
	Id                  int64     `json:"id"`
	Rank                int64     `json:"rank"`
	Name                string    `json:"name"`
	Symbol              string    `json:"symbol"`
	Slug                string    `json:"slug"`
	IsActive            int64     `json:"is_active"`
	FirstHistoricalData string    `json:"first_historical_data"`
	LastHistoricalData  string    `json:"last_historical_data"`
	Platform            *Platform `json:"platform"`
}

type Platform struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

func GetTopToken(num int64) []TokenInfo {
	url := CmcUrl

	req, _ := http.NewRequest("GET", url+"&limit="+strconv.FormatInt(num, 10), nil)

	req.Header.Add("X-CMC_PRO_API_KEY", CmcAuthVal)

	res, _ := http.DefaultClient.Do(req)
	if res == nil {
		logger.Info("GetTopToken err, ", "num", num)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	data := string(body)

	var tokenResp *TokenResp

	err := json.Unmarshal([]byte(data), &tokenResp)
	if err != nil {
		logger.Error("GetTopToken parse err ", "msg", err)
		return nil
	}

	return tokenResp.Data

}

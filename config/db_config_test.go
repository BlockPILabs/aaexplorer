package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

/*
```json
[{"Group":"default","Schema":{"AAContractInteract":"","AAHotTokenStatistic":"","AAUserOpsCalldata":"","AAUserOpsInfo":"","AaAccountData":"","AaBlockInfo":"","AaBlockSync":"","AaTransactionInfo":"","Account":"","AssetChangeTrace":"","BlockDataDecode":"","BlockScanRecord":"","BundlerInfo":"","BundlerStatisDay":"","BundlerStatisHour":"","DailyStatisticDay":"","DailyStatisticHour":"","FactoryInfo":"","FactoryStatisDay":"","FactoryStatisHour":"","FunctionSignature":"","HotAATokenStatistic":"","MevInfo":"","Network":"","PaymasterInfo":"","PaymasterStatisDay":"","PaymasterStatisHour":"","TokenPriceInfo":"","TransactionDecode":"","TransactionReceiptDecode":"","UserAssetInfo":"","UserOpTypeStatistic":"","WhaleStatisticDay":"","WhaleStatisticHour":""},"Type":"postgres","Host":"127.0.0.1","Port":5432,"User":"postgres","Pass":"root","Name":"postgres","ApplicationName":"aim","MaxIdleCount":50,"MaxOpenConns":100,"MaxLifetime":3600000000000,"Debug":true,"SslMode":""}]
```
*/
func TestJson(t *testing.T) {
	bytes, err := json.Marshal(DefaultDatabaseConfig())
	if err != nil {
		return
	}
	fmt.Println(string(bytes))
}

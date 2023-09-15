package task

import (
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/task/aa"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/shopspring/decimal"
	"sync"
)

type UserOperationEvent struct {
	OpsHash       string
	Sender        string
	Paymaster     string
	Nonce         int64
	Success       int
	ActualGasCost int64
	ActualGasUsed int64
	Target        string
	Factory       string
}

type CallDetail struct {
	target string
	value  decimal.Decimal
	data   string
	source string
}

type _evmParser struct {
	logger          log.Logger
	config          *config.Config
	startBlock      map[string]int64
	abi             abi.ABI
	handleOpsMethod *abi.Method
}

type parserBlock struct {
	block       *ent.BlockDataDecode
	transitions []*parserTransaction
	userOpInfo  *ent.AaBlockInfo
	aaAccounts  *sync.Map
}

func (b *parserBlock) AaAccountData(address string) *ent.AaAccountData {
	a, _ := b.aaAccounts.LoadOrStore(address, &ent.AaAccountData{ID: address})
	return a.(*ent.AaAccountData)
}
func (b *parserBlock) AaAccountDataSlice() ent.AaAccountDataSlice {
	s := ent.AaAccountDataSlice{}
	b.aaAccounts.Range(func(key, value any) bool {
		s = append(s, value.(*ent.AaAccountData))
		return true
	})
	return s
}

type parserTransaction struct {
	transaction     *ent.TransactionDecode
	receipt         *ent.TransactionReceiptDecode
	userOpInfo      *ent.AaTransactionInfo
	userops         ent.AAUserOpsInfos
	logs            []*aa.Log
	userOpsCalldata ent.AAUserOpsCalldataSlice
}

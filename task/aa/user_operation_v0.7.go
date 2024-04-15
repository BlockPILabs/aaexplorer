package aa

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type UserOperationV07 struct {
	Sender             common.Address `json:"sender"               mapstructure:"sender"               validate:"required"`
	Nonce              *big.Int       `json:"nonce"                mapstructure:"nonce"                validate:"required"`
	InitCode           []byte         `json:"initCode"             mapstructure:"initCode"             validate:"required"`
	CallData           []byte         `json:"callData"             mapstructure:"callData"             validate:"required"`
	PreVerificationGas *big.Int       `json:"preVerificationGas"   mapstructure:"preVerificationGas"   validate:"required"`
	PaymasterAndData   []byte         `json:"paymasterAndData"     mapstructure:"paymasterAndData"     validate:"required"`
	Signature          []byte         `json:"signature"            mapstructure:"signature"            validate:"required"`
	GasFees            []byte         `json:"gasFees"              mapstructure:"gasFees"              validate:"required"`
	AccountGasLimits   []byte         `json:"accountGasLimits"     mapstructure:"accountGasLimits"     validate:"required"`
}

func (up *UserOperationV07) UnpackGasFees() (maxFeePerGas *big.Int, maxPriorityFeePerGas *big.Int) {
	gasFees := up.GasFees
	maxFeePerGas = new(big.Int).SetBytes(gasFees[:16])
	maxPriorityFeePerGas = new(big.Int).SetBytes(gasFees[16:])
	return maxFeePerGas, maxPriorityFeePerGas
}

func (up *UserOperationV07) UnpackAccountGasLimits() (callGasLimit *big.Int, verificationGasLimit *big.Int) {
	accountGasLimits := up.AccountGasLimits
	callGasLimit = new(big.Int).SetBytes(accountGasLimits[:16])
	verificationGasLimit = new(big.Int).SetBytes(accountGasLimits[16:])
	return callGasLimit, verificationGasLimit
}

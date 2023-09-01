package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strconv"
	"strings"
)

func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func IsAddress(address string) bool {
	if len(address) != 40 {
		return false
	}
	for i := 0; i < 15; i++ {
		if address[i] != '0' {
			return true
		}
	}
	return false
}

func HexToDecimal(hexStr string) *big.Int {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	decimal := new(big.Int)
	_, success := decimal.SetString(hexStr, 16)
	if !success {
		return nil
	}

	return decimal
}

func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}

func HexToDecimalInt(hexStr string) *int {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	decimal := new(big.Int)
	_, success := decimal.SetString(hexStr, 16)
	if !success {
		return nil
	}
	res, err := strconv.Atoi(decimal.String())
	if err != nil {
		return nil
	}
	return &res
}

func HexToAddress(hexStr string) string {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	address := strings.ToLower(common.HexToAddress(hexStr).String())
	return address
}

func Substring(input string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(input) {
		end = len(input)
	}

	return input[start:end]
}

func SubstringFromIndex(input string, index int) string {
	if index < 0 || index >= len(input) {
		return ""
	}
	return input[index:]
}

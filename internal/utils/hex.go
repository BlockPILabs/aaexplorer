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

func Fix0x(input string) string {
	if Has0xPrefix(input) {
		return input
	}
	return "0x" + input
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// IsHex validates whether each byte is valid hexadecimal string.
func IsHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// IsHexSting validates whether each byte is valid hexadecimal string.
func IsHexSting(str string) bool {
	if Has0xPrefix(str) {
		str = str[2:]
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
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

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func IsHexAddress(s string) bool {
	if Has0xPrefix(s) {
		s = s[2:]
	}
	return IsAddress(s)
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

func HexToDecimalInt(hexStr string) int {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	decimal := new(big.Int)
	_, success := decimal.SetString(hexStr, 16)
	if !success {
		return 0
	}
	res, err := strconv.Atoi(decimal.String())
	if err != nil {
		return 0
	}
	return res
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

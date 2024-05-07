package util

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func Add0xPrefix(input string) string {
	if Has0xPrefix(input) {
		return input
	}
	return "0x" + input
}

// isHex validates whether each byte is valid hexadecimal string.
func IsHex(str string) bool {
	if Has0xPrefix(str) {
		str = str[2:]
	}
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

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func RecoverSignFromMetamask(dataStr, signHex string) (string, error) {

	if !IsHex(signHex) {
		return "", errors.New("sign hex error")
	}

	msg := accounts.TextHash([]byte(dataStr))

	signBytes, err := hexutil.Decode(signHex)
	if err != nil {
		return "", err
	}

	signBytes[crypto.RecoveryIDOffset] -= 27

	sigPublicKey, err := crypto.SigToPub(msg, signBytes)
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*sigPublicKey).Hex(), nil
}

func IsSameAddress(addr1, addr2 string) bool {
	addr1Hex := addr1
	if !strings.HasPrefix(addr1, "0x") {
		addr1Hex = "0x" + addr1
	}
	addr1Hex = strings.ToLower(addr1Hex)

	addr2Hex := addr2
	if !strings.HasPrefix(addr2, "0x") {
		addr2Hex = "0x" + addr2
	}
	addr2Hex = strings.ToLower(addr2Hex)

	return addr1Hex == addr2Hex
}

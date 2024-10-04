package util

import (
	"math/rand"
	"strings"
	"time"
)

const characters = "abcdefghijklmnopqrstuvwxyz0123456789"

var letters = characters[0:26]

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func byteWithCharset(charset string) byte {
	return charset[seededRand.Intn(len(charset))]
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func VariableName() string {
	var sb strings.Builder
	length := seededRand.Intn(12)
	sb.WriteByte(byteWithCharset(letters))
	sb.WriteString(stringWithCharset(length, characters))
	return sb.String()
}

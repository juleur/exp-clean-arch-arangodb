package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoding(t *testing.T) {
	nbsStr := []string{"1243", "8453", "3454930", "9453421"}
	for _, nbStr := range nbsStr {
		encodedString, _ := EncodeString(nbStr)
		t.Log(encodedString)
	}
}

func TestDecoding(t *testing.T) {
	assert := assert.New(t)

	encodedStrs := []string{"YCj2a", "176hSn", "1eb09O2fUs", "1mZnnYG1tv"}
	expectedStrs := []string{"1243", "8453", "3454930", "9453421"}

	for i, encStr := range encodedStrs {
		decodedStr, _ := DecodeString(encStr)
		assert.Equal(expectedStrs[i], decodedStr)
	}
}

package utils

import (
	"errors"

	"github.com/eknkc/basex"
)

const EncodingErr = "Encoding error"
const DecodingErr = "Decoding error"

// le ‘u‘ a été retiré puisqu'il servira de séparateur avec le session ID
const BASE = "0123456789abcdefghijklmnopqrstvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func EncodeString(s string) (string, error) {
	enc, err := basex.NewEncoding(BASE)
	if err != nil {
		return "", errors.New(EncodingErr)
	}
	encodedNumber := enc.Encode([]byte(s))

	return encodedNumber, nil
}

func DecodeString(s string) (string, error) {
	enc, err := basex.NewEncoding(BASE)
	if err != nil {
		return "", errors.New(DecodingErr)
	}

	numberStr, err := enc.Decode(s)
	if err != nil {
		return "", errors.New(DecodingErr)
	}

	return string(numberStr), nil
}

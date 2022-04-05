package crypt

import (
	"encoding/base64"
	"errors"
	"strconv"
)

func FctHashSecret(_b string, _salt string) (string, error) {

	// Validate the salt
	var minor rune
	var offset int
	if _salt[0] != '$' || _salt[1] != '2' {
		return "", errors.New("Invalid salt version: " + _salt[0:2])
	}
	if _salt[2] == '$' {
		minor = 0
		offset = 3
	} else {
		minor = rune(_salt[2])
		if (minor != 'a' && minor != 'b' && minor != 'y') || _salt[3] != '$' {
			return "", errors.New("Invalid salt version: " + _salt[2:4])
		}
		offset = 4
	}

	// Extract number of rounds
	if _salt[offset+2] > '$' {
		return "", errors.New("missing salt rounds")
	}

	r1, err := strconv.Atoi(_salt[offset : offset+1])
	if err != nil {
		return "", errors.New("invalid r1 decimal value")
	}
	r1 *= 10

	r2, err := strconv.Atoi(_salt[offset+1 : offset+2])
	if err != nil {
		return "", errors.New("invalid r2 decimal value")
	}

	rounds := r1 + r2
	realSalt := _salt[offset+3 : offset+25]

	/*
		var s string
		if minor >= 'a' {
			s = "\x00"
		}
		passwordb := []byte(s)
	*/

	saltb, err := base64Decode([]byte(realSalt))
	if err != nil {
		return "", errors.New("can't decode string to base64")
	}

	var res string
	res += "$2"
	if minor >= 'a' {
		res += string(minor)
	}
	res += "$"
	if rounds < 10 {
		res += "0"
	}
	res += strconv.Itoa(rounds)
	res += "$"
	res += string(base64Encode(saltb))
	res += string(base64Encode([]byte(_b)))

	return res, nil
}

const alphabet = "./ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var bcEncoding = base64.NewEncoding(alphabet)

func base64Encode(src []byte) []byte {
	n := bcEncoding.EncodedLen(len(src))
	dst := make([]byte, n)
	bcEncoding.Encode(dst, src)
	for dst[n-1] == '=' {
		n--
	}
	return dst[:n]
}

func base64Decode(src []byte) ([]byte, error) {
	numOfEquals := 4 - (len(src) % 4)
	for i := 0; i < numOfEquals; i++ {
		src = append(src, '=')
	}

	dst := make([]byte, bcEncoding.DecodedLen(len(src)))
	n, err := bcEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

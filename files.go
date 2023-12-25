package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func FileMD5Hash(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	h := md5.Sum(b)
	res := hex.EncodeToString(h[:])
	return res, nil
}

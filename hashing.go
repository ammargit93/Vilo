package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func GenerateCommitHash(commitMsg string, stagedFiles []string) string {
	data := commitMsg + "|" + time.Now().UTC().String()
	for _, file := range stagedFiles {
		data += "|" + file
	}
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

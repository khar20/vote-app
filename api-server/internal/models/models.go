package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// Vote represents the structure of a vote
type Vote struct {
	Timestamp string   `json:"timestamp"`
	Vote      []string `json:"vote"`
}

type VoteBlock struct {
	Timestamp string   `json:"timestamp"`
	Vote      []string `json:"vote"`
	PrevHash  string   `json:"prev_hash"`
	Hash      string   `json:"hash"`
}

// ComputeHash generates a SHA-256 hash for a vote
func (v *VoteBlock) ComputeHash() string {
	data, _ := json.Marshal(v)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

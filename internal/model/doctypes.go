package model

import "time"

const (
	ItemTypeLoginPass = iota
	ItemTypeText
	ItemTypeBinary
	ItemTypeCard
)

type Metainfo struct {
	Date time.Time `json:"creationDate"`
	Text string    `json:"info"`
}

type Pair struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PairInfo struct {
	Pair
	Metainfo `json:"meta"`
}

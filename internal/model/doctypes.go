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

type (
	LimitCtxKey  struct{}
	OffsetCtxKey struct{}
)

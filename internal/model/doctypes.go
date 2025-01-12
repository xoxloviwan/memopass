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

type File struct {
	Name string `json:"name"`
	Blob []byte `json:"blob"`
}

type FileInfo struct {
	File
	Metainfo `json:"meta"`
}

type Card struct {
	Number   string `json:"ccn"`
	Exp      string `json:"exp"`
	VerifVal string `json:"cvv"`
}

type CardInfo struct {
	Card
	Metainfo `json:"meta"`
}

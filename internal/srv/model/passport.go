package model

import "time"

type Metainfo struct {
	Date time.Time `json:"creationDate"`
	Text string    `json:"info"`
}

type Pairs struct {
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Meta     Metainfo `json:"meta"`
}

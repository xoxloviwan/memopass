package model

type Creds struct {
	User string `json:"login"`
	Pwd  string `json:"password"`
}

type User struct {
	ID   int
	Name string
	Hash []byte
}

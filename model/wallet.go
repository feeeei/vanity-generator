package model

type Wallet struct {
	Private string `json:"private_key"`
	Public  string `json:"address"`
}

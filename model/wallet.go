package model

type Wallet struct {
	Private  string `json:"private_key,omitempty"`
	Public   string `json:"address"`
	Mnemonic string `json:"mnemonic,omitempty"`
	Path     string `json:"path,omitempty"`
}

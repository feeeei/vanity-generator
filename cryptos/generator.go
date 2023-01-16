package cryptos

import "vanity-generator/model"

type Generator interface {
	Difficulty(prefix, suffix string) (count int64)
	DoSingle(prefix, suffix string) *model.Wallet
}

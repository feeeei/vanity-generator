package cryptos

import "vanity-generator/model"

// DefaultMnemonicWords 是助记词的默认单词数：MetaMask、TronLink、polkadot.js 都同时支持 12 和 24 词，默认取更长的 24 词
const DefaultMnemonicWords = 24

type Generator interface {
	Difficulty(prefix, suffix string) (count int64)
	DoSingle(prefix, suffix string) *model.Wallet
}

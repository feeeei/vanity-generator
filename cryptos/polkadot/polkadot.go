package polkadot

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/tyler-smith/go-bip39"
	"math"
	"strings"
	"vanity-generator/cryptos"
	"vanity-generator/model"
)

var _ cryptos.Generator = (*PolkadotGenerator)(nil)

type PolkadotGenerator struct {
}

func NewPolkadotGenerator() *PolkadotGenerator {
	return &PolkadotGenerator{}
}

func (g *PolkadotGenerator) Difficulty(prefix string, suffix string) (count int64) {
	if len(prefix) > 0 {
		prefix = prefix[1:]
	}
	difficulty := 1.0
	if len(prefix) > 0 {
		difficulty = firstDiff(prefix[0])
	}
	if len(prefix) > 1 {
		difficulty *= math.Pow(58, float64(len(prefix)-1))
	}

	return int64(difficulty * math.Pow(58, float64(len(suffix))))
}

func (g *PolkadotGenerator) DoSingle(prefix string, suffix string) *model.Wallet {
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	keyringPair, _ := signature.KeyringPairFromSecret(mnemonic, 0)
	if strings.HasPrefix(keyringPair.Address, prefix) && strings.HasSuffix(keyringPair.Address, suffix) {
		return &model.Wallet{
			Private: mnemonic,
			Public:  keyringPair.Address,
		}
	}
	return nil
}

func firstDiff(char byte) float64 {
	switch char {
	case 2, 3, 4, 5, 6:
		return 5.988
	default:
		return 297.78
	}
}

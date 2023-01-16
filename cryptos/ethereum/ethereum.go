package ethereum

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math"
	"strings"
	"vanity-generator/cryptos"
	"vanity-generator/model"
)

var _ cryptos.Generator = (*EthereumGenerator)(nil)

type EthereumGenerator struct {
}

func NewEthereumGenerator() *EthereumGenerator {
	return &EthereumGenerator{}
}

func (g *EthereumGenerator) Difficulty(prefix, suffix string) (count int64) {
	if len(prefix) > 2 {
		prefix = prefix[2:]
	}
	num, letter := word(prefix + suffix)
	return int64(math.Pow(16, float64(num)) * math.Pow(32, float64(letter)))
}

func (g *EthereumGenerator) DoSingle(prefix, suffix string) *model.Wallet {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	public := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	if strings.HasPrefix(public, prefix) && strings.HasSuffix(public, suffix) {
		privateKeyBytes := crypto.FromECDSA(privateKey)
		private := hexutil.Encode(privateKeyBytes)[2:]
		return &model.Wallet{
			Private: private,
			Public:  public,
		}
	}
	return nil
}

// word 计算str中的数字、字母字数
func word(str string) (num, letter int) {
	for _, char := range str {
		if char >= '0' && char <= '9' {
			num++
			continue
		}
		if char >= 'a' && char <= 'z' {
			letter++
			continue
		}
		if char >= 'A' && char <= 'Z' {
			letter++
			continue
		}
	}
	return
}

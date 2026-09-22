package ethereum

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math"
	"strings"
	"vanity-generator/cryptos"
	"vanity-generator/cryptos/hdwallet"
	"vanity-generator/model"
)

var _ cryptos.Generator = (*EthereumGenerator)(nil)

type EthereumGenerator struct {
	words int // 助记词单词数，0 表示直接生成随机私钥
}

// NewEthereumGenerator 构造生成器，words 为 0 时直接随机私钥，为 12 或 24 时按 BIP39 助记词经 m/44'/60'/0'/0/0 派生
func NewEthereumGenerator(words int) *EthereumGenerator {
	return &EthereumGenerator{words: words}
}

func (g *EthereumGenerator) Difficulty(prefix, suffix string) (count int64) {
	if len(prefix) < 2 {
		prefix = ""
	} else {
		prefix = prefix[2:]
	}
	num, letter := word(prefix + suffix)
	return int64(math.Pow(16, float64(num)) * math.Pow(32, float64(letter)))
}

func (g *EthereumGenerator) DoSingle(prefix, suffix string) *model.Wallet {
	mnemonic, privateKey, err := g.newKey()
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
		wallet := &model.Wallet{
			Private: private,
			Public:  public,
		}
		if mnemonic != "" {
			wallet.Mnemonic = mnemonic
			wallet.Path = hdwallet.Path(hdwallet.CoinEthereum)
		}
		return wallet
	}
	return nil
}

// newKey 生成一把新私钥，助记词模式下同时返回助记词
func (g *EthereumGenerator) newKey() (mnemonic string, key *ecdsa.PrivateKey, err error) {
	if g.words == 0 {
		key, err = crypto.GenerateKey()
		return "", key, err
	}
	return hdwallet.NewKey(g.words, hdwallet.CoinEthereum)
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

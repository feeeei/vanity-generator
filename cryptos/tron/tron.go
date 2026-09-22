package tron

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"log"
	"math"
	"strings"
	"sync"
	"vanity-generator/common"
	"vanity-generator/cryptos"
	"vanity-generator/cryptos/hdwallet"
	"vanity-generator/model"
)

var _ cryptos.Generator = (*TronGenerator)(nil)

type TronGenerator struct {
	addrPool *sync.Pool
	words    int // 助记词单词数，0 表示直接生成随机私钥
}

// NewTronGenerator 构造生成器，words 为 0 时直接随机私钥，为 12 或 24 时按 BIP39 助记词经 m/44'/195'/0'/0/0 派生
func NewTronGenerator(words int) *TronGenerator {
	return &TronGenerator{
		words: words,
		addrPool: &sync.Pool{
			New: func() any {
				addr := [25]byte{0x41}
				return addr[:]
			},
		},
	}
}

func (g *TronGenerator) Difficulty(prefix, suffix string) int64 {
	if len(prefix) > 0 {
		prefix = prefix[1:]
	}
	// 第一位奇葩规则，只有A-H J-N P-Z 9 这25个字符
	// 其中 A-H、J-N、P-Y 的概率为平均的，出现概率为4.285036%
	// Z 的概率为 1.3126%
	// 9 的概率为 0.1315%
	// 后续字符出现概率均等
	first := 1.0
	if len(prefix) > 0 {
		first = firstDiff(prefix[0])
		prefix = prefix[1:]
	}

	diff := first * math.Pow(58, float64(len(prefix)+len(suffix)))
	return int64(diff)
}

func (g *TronGenerator) DoSingle(prefix, suffix string) *model.Wallet {
	mnemonic, privateKey, err := g.newKey()
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)

	hex := crypto.PubkeyToAddress(*publicKeyECDSA).Bytes()
	public := g.hexToBase58CheckAddress(hex)

	if strings.HasPrefix(public, prefix) && strings.HasSuffix(public, suffix) {
		wallet := &model.Wallet{
			Private: hexutil.Encode(crypto.FromECDSA(privateKey))[2:],
			Public:  public,
		}
		if mnemonic != "" {
			wallet.Mnemonic = mnemonic
			wallet.Path = hdwallet.Path(hdwallet.CoinTron)
		}
		return wallet
	}
	return nil
}

// newKey 生成一把新私钥，助记词模式下同时返回助记词
func (g *TronGenerator) newKey() (mnemonic string, key *ecdsa.PrivateKey, err error) {
	if g.words == 0 {
		key, err = crypto.GenerateKey()
		return "", key, err
	}
	return hdwallet.NewKey(g.words, hdwallet.CoinTron)
}

func (g *TronGenerator) hexToBase58CheckAddress(hex []byte) string {
	addr := g.addrPool.Get().([]byte)
	defer g.addrPool.Put(addr)

	copy(addr[1:], hex)
	sum := common.S256(common.S256(addr[:21]))[:4]
	copy(addr[21:], sum[:4])
	return base58.Encode(addr)
}

func firstDiff(char byte) float64 {
	if char >= 'A' && char < 'Y' {
		return 23.337 // 100/4.285036=23.337
	}
	if char == 'Z' {
		return 76.185 // 100/1.3126=76.185
	}
	if char == '9' {
		return 760.456 // 100/0.1315=760.456
	}
	return 1
}

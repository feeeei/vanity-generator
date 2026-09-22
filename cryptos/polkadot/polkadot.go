package polkadot

import (
	"crypto/rand"
	"crypto/sha512"
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/tyler-smith/go-bip39"
	"github.com/vedhavyas/go-subkey/v2"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"math"
	"os"
	"strings"
	"sync"
	"unsafe"
	"vanity-generator/cryptos"
	"vanity-generator/model"
)

var _ cryptos.Generator = (*PolkadotGenerator)(nil)

type PolkadotGenerator struct {
	writer      io.Writer
	entropyPool *sync.Pool
}

// NewPolkadotGenerator 构造生成器，words 为助记词单词数（12 或 24），Substrate 钱包始终由 BIP39 助记词派生
func NewPolkadotGenerator(words int) *PolkadotGenerator {
	out, _ := os.Open(os.DevNull)
	entropySize := words / 3 * 4 // 12 词 16 字节，24 词 32 字节
	return &PolkadotGenerator{
		writer: out,
		entropyPool: &sync.Pool{
			New: func() any {
				return make([]byte, entropySize)
			},
		},
	}
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
	entropy := g.entropyPool.Get().([]byte)
	defer g.entropyPool.Put(entropy)

	_, _ = rand.Read(entropy)

	address := g.address(entropy)

	if strings.HasPrefix(address, prefix) && strings.HasSuffix(address, suffix) {
		mnemonic, _ := bip39.NewMnemonic(entropy)
		return &model.Wallet{
			Mnemonic: mnemonic,
			Public:   address,
		}
	}
	return nil
}

// address 按 Substrate 规则从 BIP39 熵推导 sr25519 密钥并编码为 SS58 地址：
// mini secret 取 PBKDF2(熵, "mnemonic", 2048 轮) 的前 32 字节，与 polkadot.js 一致
func (g *PolkadotGenerator) address(entropy []byte) string {
	seed := pbkdf2.Key(entropy, []byte("mnemonic"), 2048, 64, sha512.New)
	ms, _ := schnorrkel.NewMiniSecretKeyFromRaw(*(*[32]byte)(unsafe.Pointer(&seed[0])))

	code := ms.Public().Encode()
	return subkey.SS58Encode(code[:], 0)
}

// firstDiff 是地址第二位字符的期望尝试次数。主网地址以 1 开头，其后是 34 字节大端整数的 base58 编码：
// 该整数多数情况下有 47 位，首位只能是 1 到 5，即字符 2 到 6，其中 6 只占部分区间；
// 约 17% 的情况整数只有 46 位，此时第二位在 58 个字符里近似均匀。数值按 60 万个随机 sr25519 公钥统计得出
func firstDiff(char byte) float64 {
	switch char {
	case '2', '3', '4', '5':
		return 5.7 // 出现概率各约 17.5%
	case '6':
		return 7.2 // 出现概率约 13.8%
	default:
		return 330 // 其余字符出现概率各约 0.3%
	}
}

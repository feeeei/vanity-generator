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

func NewPolkadotGenerator() *PolkadotGenerator {
	out, _ := os.Open(os.DevNull)
	return &PolkadotGenerator{
		writer: out,
		entropyPool: &sync.Pool{
			New: func() any {
				addr := [128 / 8]byte{}
				return addr[:]
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

	seed := pbkdf2.Key(entropy, []byte("mnemonic"), 2048, 64, sha512.New)
	ms, _ := schnorrkel.NewMiniSecretKeyFromRaw(*(*[32]byte)(unsafe.Pointer(&seed[0])))

	code := ms.Public().Encode()
	address := subkey.SS58Encode(code[:], 0)

	if strings.HasPrefix(address, prefix) && strings.HasSuffix(address, suffix) {
		mnemonic, _ := bip39.NewMnemonic(entropy)
		return &model.Wallet{
			Private: mnemonic,
			Public:  address,
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

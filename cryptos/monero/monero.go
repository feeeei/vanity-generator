package monero

import (
	"crypto/rand"
	"log"
	"math"
	"strings"
	"vanity-generator/cryptos"
	"vanity-generator/model"
)

var _ cryptos.Generator = (*MoneroGenerator)(nil)

// halfAddressPrefixLength 是 HalfAddress 与完整地址一致的前缀长度：base58 按 8 字节分块独立编码，
// 前 4 块（网络字节 + spend 公钥前 31 字节）共 44 个字符，与 view 公钥无关
const halfAddressPrefixLength = 44

type MoneroGenerator struct {
}

func NewMoneroGenerator() *MoneroGenerator {
	return &MoneroGenerator{}
}

// Difficulty 计算期望尝试次数。主网地址固定以 4 开头；第二位由前 8 字节（0x12 加 spend 公钥前 7 字节）决定，
// 取值区间 [0x1200000000000000, 0x12ffffffffffffff] 只覆盖 1-9、A、B，其中 1 和 B 只占各自 base58 位的一部分；
// 其后每一位近似均匀分布
func (g *MoneroGenerator) Difficulty(prefix, suffix string) int64 {
	if len(prefix) > 0 {
		prefix = prefix[1:]
	}
	diff := 1.0
	if len(prefix) > 0 {
		diff = secondDiff(prefix[0])
		prefix = prefix[1:]
	}
	return int64(diff * math.Pow(58, float64(len(prefix)+len(suffix))))
}

func secondDiff(char byte) float64 {
	switch char {
	case '1':
		return 25.702 // 出现概率 3.891%
	case 'B':
		return 29.961 // 出现概率 3.338%
	default:
		return 9.701 // 2-9、A 出现概率各 10.308%
	}
}

func (g *MoneroGenerator) DoSingle(prefix, suffix string) *model.Wallet {
	var seed [32]byte
	if _, err := rand.Read(seed[:]); err != nil {
		log.Fatal(err)
	}

	key := HalfKeyFromSeed(&seed)
	// 前缀只由 spend 公钥决定，先用半地址过滤，不匹配的候选省掉 view 密钥和校验和的计算
	if len(prefix) <= halfAddressPrefixLength {
		half := key.HalfAddress(MoneroMainNetwork)
		if !strings.HasPrefix(half[:halfAddressPrefixLength], prefix) {
			return nil
		}
	}

	key.HalfToFull()
	address := key.Address(MoneroMainNetwork)
	if strings.HasPrefix(address, prefix) && strings.HasSuffix(address, suffix) {
		return &model.Wallet{
			Mnemonic: strings.Join(English.Encode(key.Seed())[:], " "),
			Public:   address,
		}
	}
	return nil
}

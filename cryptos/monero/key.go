// 本文件移植自 ekyu.moe/vanity-monero 的 key.go、wrapper.go 和 variant.go
// Copyright (c) 2018 Equim，MIT 许可，完整许可文本见同目录 THIRD_PARTY_LICENSES
// 改动：曲线运算由 moneroutil 改为 filippo.io/edwards25519，keccak 实例复用 sync.Pool，去掉了 split-key 相关方法

package monero

import (
	"hash"
	"sync"

	"filippo.io/edwards25519"
	"golang.org/x/crypto/sha3"
)

// Network 是地址开头的网络字节
type Network []byte

var (
	MoneroMainNetwork = Network{0x12}
	MoneroTestNetwork = Network{0x35}
)

// Key 由 32 字节的 ed25519 spend 私钥和 view 私钥组成。
// 为了让助记词能还原钱包，view 私钥由 spend 私钥经 Keccak256 推导，因此 spend 私钥即种子
type Key struct {
	SpendKey, ViewKey *[32]byte
}

// KeyFromSeed 用给定种子构造完整的 Key
func KeyFromSeed(seed *[32]byte) *Key {
	k := HalfKeyFromSeed(seed)
	k.HalfToFull()
	return k
}

// HalfKeyFromSeed 只推导 spend 私钥，用于靓号生成加速：地址前缀只由 spend 公钥决定，
// 不匹配的候选不必再推导 view 密钥
func HalfKeyFromSeed(seed *[32]byte) *Key {
	k := &Key{new([32]byte), new([32]byte)}
	copy(k.SpendKey[:], seed[:])
	scReduce32(k.SpendKey)
	return k
}

// HalfToFull 把 HalfKeyFromSeed 构造的 Key 补全，即从 spend 私钥推导 view 私钥
func (k *Key) HalfToFull() {
	k.ViewKey = keccak256(k.SpendKey[:])
	scReduce32(k.ViewKey)
}

func (k *Key) PublicSpendKey() *[32]byte {
	return publicKeyFromPrivateKey(k.SpendKey)
}

func (k *Key) PublicViewKey() *[32]byte {
	return publicKeyFromPrivateKey(k.ViewKey)
}

// Seed 返回助记词种子，即 spend 私钥
func (k *Key) Seed() *[32]byte {
	return k.SpendKey
}

// HalfAddress 只编码网络字节和 spend 公钥，前 44 个字符与完整地址一致，仅用于前缀匹配
func (k *Key) HalfAddress(network Network) string {
	spendPub := k.PublicSpendKey()
	return encodeBase58(network, spendPub[:])
}

// Address 把网络字节、spend 公钥、view 公钥和 Keccak256 校验和前 4 字节编码为 base58 地址
func (k *Key) Address(network Network) string {
	spendPub := k.PublicSpendKey()
	viewPub := k.PublicViewKey()
	hash := keccak256(network, spendPub[:], viewPub[:])
	return encodeBase58(network, spendPub[:], viewPub[:], hash[:4])
}

var keccakPool = sync.Pool{
	New: func() any {
		return sha3.NewLegacyKeccak256()
	},
}

func keccak256(data ...[]byte) *[32]byte {
	h := keccakPool.Get().(hash.Hash)
	for _, v := range data {
		h.Write(v)
	}
	sum := new([32]byte)
	h.Sum(sum[:0])
	h.Reset()
	keccakPool.Put(h)
	return sum
}

// scReduce32 把 32 字节小端整数对曲线阶 l 取模，等价于 Monero 的 sc_reduce32
func scReduce32(b *[32]byte) {
	var wide [64]byte
	copy(wide[:32], b[:])
	s, err := new(edwards25519.Scalar).SetUniformBytes(wide[:])
	if err != nil {
		panic(err) // wide 固定 64 字节，不会走到这里
	}
	copy(b[:], s.Bytes())
}

// publicKeyFromPrivateKey 计算 priv * G，priv 必须已经过 scReduce32
func publicKeyFromPrivateKey(priv *[32]byte) *[32]byte {
	s, err := new(edwards25519.Scalar).SetCanonicalBytes(priv[:])
	if err != nil {
		panic(err) // 私钥已规约到 [0, l)，不会走到这里
	}
	pub := new([32]byte)
	copy(pub[:], new(edwards25519.Point).ScalarBaseMult(s).Bytes())
	return pub
}

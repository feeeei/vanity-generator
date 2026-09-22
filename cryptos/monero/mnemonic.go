// 本文件移植自 ekyu.moe/vanity-monero 的 mnemonic/dict.go
// Copyright (c) 2018 Equim，MIT 许可，完整许可文本见同目录 THIRD_PARTY_LICENSES
// 改动：只保留英文词表，去掉了解码用的反查表

package monero

import (
	"encoding/binary"
	"hash/crc32"
)

const dictSize = 1626

// English 是 Monero 官方英文助记词词表，唯一前缀长度为 3
var English = newDict(englishTable, 3)

// Dict 是助记词词表
type Dict struct {
	UniquePrefixLength int
	Table              *[dictSize]string
}

func newDict(table *[dictSize]string, prefixLen int) *Dict {
	return &Dict{prefixLen, table}
}

// Encode 把 32 字节种子编码为 25 个助记词：每 4 字节（小端）编码为 3 个词，第 25 个是校验词
func (d *Dict) Encode(key *[32]byte) *[25]string {
	w := new([25]string)
	for i := 0; i < 32; i += 4 {
		x := binary.LittleEndian.Uint32(key[i : i+4])
		w1 := x % dictSize
		w2 := (x/dictSize + w1) % dictSize
		w3 := (x/dictSize/dictSize + w2) % dictSize
		w[i/4*3] = d.Table[w1]
		w[i/4*3+1] = d.Table[w2]
		w[i/4*3+2] = d.Table[w3]
	}
	w[24] = d.getChecksumWord(w)

	return w
}

// getChecksumWord 取前 24 个词各自的唯一前缀拼接后做 CRC32，对 24 取模得到校验词下标
func (d *Dict) getChecksumWord(w *[25]string) string {
	h := crc32.NewIEEE()
	for _, v := range w[:24] {
		r := string([]rune(v)[:d.UniquePrefixLength])
		h.Write([]byte(r))
	}
	sum := h.Sum32()
	idx := sum % 24

	return w[idx]
}

// 本文件移植自 github.com/paxosglobal/moneroutil 的 base58.go，只保留编码部分
// Copyright (c) 2017 Paxos，MIT 许可，完整许可文本见同目录 THIRD_PARTY_LICENSES
// 改动：每块最多 8 字节，改用 uint64 代替 big.Int 运算；尾块编码长度按 Monero 官方 encoded_block_sizes 表取值

package monero

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

const fullBlockSize = 8

// encodedBlockSizes[n] 是 n 字节的块编码后的字符数
var encodedBlockSizes = [fullBlockSize + 1]int{0, 2, 3, 5, 6, 7, 9, 10, 11}

// encodeChunk 把最多 8 字节的块按大端整数编码到 dst，高位不足以 1 补齐
func encodeChunk(dst []byte, raw []byte) {
	var num uint64
	for _, b := range raw {
		num = num<<8 | uint64(b)
	}
	for i := len(dst) - 1; i >= 0; i-- {
		dst[i] = base58Alphabet[num%58]
		num /= 58
	}
}

// encodeBase58 是 Monero 的 base58 编码：每 8 字节一块编码为 11 个字符，尾块按 encodedBlockSizes 取长度
func encodeBase58(data ...[]byte) string {
	var combined []byte
	for _, item := range data {
		combined = append(combined, item...)
	}

	rounds, rest := len(combined)/fullBlockSize, len(combined)%fullBlockSize
	fullSize := encodedBlockSizes[fullBlockSize]
	encoded := make([]byte, rounds*fullSize+encodedBlockSizes[rest])
	for i := 0; i < rounds; i++ {
		encodeChunk(encoded[i*fullSize:(i+1)*fullSize], combined[i*fullBlockSize:(i+1)*fullBlockSize])
	}
	if rest > 0 {
		encodeChunk(encoded[rounds*fullSize:], combined[rounds*fullBlockSize:])
	}
	return string(encoded)
}

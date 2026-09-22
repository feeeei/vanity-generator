package polkadot

import (
	"strings"
	"testing"

	"github.com/tyler-smith/go-bip39"
)

// 向量由 polkadot.js 的 mnemonicToMiniSecret 与 sr25519PairFromSeed 生成，SS58 前缀 0
func TestAddress(t *testing.T) {
	tests := []struct{ mnemonic, address string }{
		{"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about", "13KVd4f2a4S5pLp4gTTFezyXdPWx27vQ9vS6xBXJ9yWVd7xo"},
		{"abandon amount liar amount expire adjust cage candy arch gather drum buyer", "1bpeaHn4GX1afW8Jr3XQareAy8Tz3gpjfETfFfZT9BY2hwg"},
		{"abandon amount liar amount expire adjust cage candy arch gather drum bullet absurd math era live bid rhythm alien crouch range attend journey unaware", "12gSWHjzeFyd2aSLVCVjByvQ8DVXSVWzQV4jAdjpdzCjExyg"},
	}
	g := NewPolkadotGenerator(12)
	for _, tt := range tests {
		entropy, err := bip39.EntropyFromMnemonic(tt.mnemonic)
		if err != nil {
			t.Fatal(err)
		}
		if got := g.address(entropy); got != tt.address {
			t.Errorf("address = %s, want %s", got, tt.address)
		}
	}
}

// 主网地址均以 1 开头，前缀 1 必然匹配
func TestDoSingleWords(t *testing.T) {
	for _, words := range []int{12, 24} {
		g := NewPolkadotGenerator(words)
		wallet := g.DoSingle("1", "")
		if wallet == nil {
			t.Fatal("expected a wallet")
		}
		if n := len(strings.Fields(wallet.Mnemonic)); n != words {
			t.Errorf("mnemonic has %d words, want %d", n, words)
		}
		if wallet.Private != "" {
			t.Errorf("polkadot wallet should only output the mnemonic: %+v", wallet)
		}
		entropy, err := bip39.EntropyFromMnemonic(wallet.Mnemonic)
		if err != nil {
			t.Fatal(err)
		}
		if got := g.address(entropy); got != wallet.Public {
			t.Errorf("address %s does not match the mnemonic (%s)", wallet.Public, got)
		}
	}
}

func TestDifficulty(t *testing.T) {
	g := NewPolkadotGenerator(12)
	tests := []struct {
		prefix, suffix string
		want           int64
	}{
		{"", "", 1},
		{"1", "", 1},
		{"12", "", 5},
		{"16", "", 7},
		{"1a", "", 330},
		{"12x", "", 330}, // 5.7 * 58 = 330.6
		{"", "xy", 58 * 58},
	}
	for _, tt := range tests {
		if got := g.Difficulty(tt.prefix, tt.suffix); got != tt.want {
			t.Errorf("Difficulty(%q, %q) = %d, want %d", tt.prefix, tt.suffix, got, tt.want)
		}
	}
}

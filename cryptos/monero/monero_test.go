package monero

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"
)

func mustSeed(t *testing.T, hexSeed string) *[32]byte {
	t.Helper()
	b, err := hex.DecodeString(hexSeed)
	if err != nil || len(b) != 32 {
		t.Fatalf("bad seed %q: %v", hexSeed, err)
	}
	seed := new([32]byte)
	copy(seed[:], b)
	return seed
}

// 第一组向量来自 monero-python 的公开测试用例，第二组是用 monero-python 对全 0xff 种子推导的结果
func TestKeyFromSeed(t *testing.T) {
	tests := []struct {
		name, seed, spendKey, viewKey, address, mnemonic string
	}{
		{
			name:     "monero-python 公开向量",
			seed:     "8ffa9f586b86d294d93731765d192765311bddc76a4fa60311f8af36bbf6fb06",
			spendKey: "8ffa9f586b86d294d93731765d192765311bddc76a4fa60311f8af36bbf6fb06",
			viewKey:  "be300e27b32e987458f5d518e0f360f62300e9132ee44df6887eda58051cba0e",
			address:  "48MjqqDYoxo4BkFtRX2hqS3TjvAjBRygX8FKv7ovEp7QW8e6TwwoYr7DyGtXYGSzbW9nWKuyhyafxU7pBQDxL16fPQvewNL",
			mnemonic: "wedge going quick racetrack auburn physics lectures light waist axes whipped habitat square awkward together injury niece nugget guarded hive obnoxious waxing faked folding square",
		},
		{
			name:     "全 0xff 种子，需要 sc_reduce32 规约",
			seed:     strings.Repeat("ff", 32),
			spendKey: "1c95988d7431ecd670cf7d73f45befc6feffffffffffffffffffffffffffff0f",
			viewKey:  "9fe83aa6104612b587eb2e6ee1f0c929f85ce047804a789f4d579f9d2e20de0b",
			address:  "49voQEbjouUQSDikRWKUt1PGbS47TBde4hiGyftN46CvTDd8LXCaimjHRGtofCJwY5Ed5QhYwc12P15AH5w7SxUAMCz1nr1",
			mnemonic: "powder directed sayings enmity bacon vapidly entrance bumper noodles iguana sleepless nasty flying soil software foamy solved soggy foamy solved soggy jury yawning ankle solved",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := KeyFromSeed(mustSeed(t, tt.seed))
			if got := hex.EncodeToString(key.SpendKey[:]); got != tt.spendKey {
				t.Errorf("spend key = %s, want %s", got, tt.spendKey)
			}
			if got := hex.EncodeToString(key.ViewKey[:]); got != tt.viewKey {
				t.Errorf("view key = %s, want %s", got, tt.viewKey)
			}
			if got := key.Address(MoneroMainNetwork); got != tt.address {
				t.Errorf("address = %s, want %s", got, tt.address)
			}
			if got := strings.Join(English.Encode(key.Seed())[:], " "); got != tt.mnemonic {
				t.Errorf("mnemonic = %s, want %s", got, tt.mnemonic)
			}
		})
	}
}

// 半地址的前 halfAddressPrefixLength 个字符必须与完整地址一致，主网地址固定 95 个字符且以 4 开头
func TestHalfAddressPrefix(t *testing.T) {
	for i := 0; i < 200; i++ {
		seed := new([32]byte)
		if _, err := rand.Read(seed[:]); err != nil {
			t.Fatal(err)
		}
		key := HalfKeyFromSeed(seed)
		half := key.HalfAddress(MoneroMainNetwork)
		key.HalfToFull()
		full := key.Address(MoneroMainNetwork)

		if len(full) != 95 || full[0] != '4' {
			t.Fatalf("unexpected address %s", full)
		}
		if half[:halfAddressPrefixLength] != full[:halfAddressPrefixLength] {
			t.Fatalf("half address %s is not a prefix of %s", half, full)
		}
	}
}

func TestDifficulty(t *testing.T) {
	g := NewMoneroGenerator()
	tests := []struct {
		prefix, suffix string
		want           int64
	}{
		{"", "", 1},
		{"4", "", 1},
		{"42", "", 9},
		{"41", "", 25},
		{"4B", "", 29},
		{"42abc", "", 1892781}, // 9.701 * 58^3 = 1892781.5
		{"", "xyz", 58 * 58 * 58},
		{"4A", "z", 562}, // 9.701 * 58 = 562.7
	}
	for _, tt := range tests {
		if got := g.Difficulty(tt.prefix, tt.suffix); got != tt.want {
			t.Errorf("Difficulty(%q, %q) = %d, want %d", tt.prefix, tt.suffix, got, tt.want)
		}
	}
}

// 主网地址均以 4 开头，前缀 4 必然匹配
func TestDoSingle(t *testing.T) {
	wallet := NewMoneroGenerator().DoSingle("4", "")
	if wallet == nil {
		t.Fatal("expected a wallet")
	}
	if n := len(strings.Fields(wallet.Mnemonic)); n != 25 {
		t.Errorf("mnemonic has %d words, want 25", n)
	}
	if wallet.Private != "" || wallet.Path != "" {
		t.Errorf("monero wallet should only output the mnemonic: %+v", wallet)
	}
	if len(wallet.Public) != 95 || wallet.Public[0] != '4' {
		t.Errorf("unexpected address %s", wallet.Public)
	}
}

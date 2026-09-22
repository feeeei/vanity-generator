package tron

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"vanity-generator/cryptos/hdwallet"
)

// 向量由 tronweb 生成，路径 m/44'/195'/0'/0/0
func TestMnemonicAddress(t *testing.T) {
	tests := []struct{ mnemonic, address string }{
		{"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about", "TUEZSdKsoDHQMeZwihtdoBiN46zxhGWYdH"},
		{"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art", "TEfhiqsW1SdN44DeHrAWVmbyr8ZbvChrtS"},
		{"abandon amount liar amount expire adjust cage candy arch gather drum buyer", "TSRqUBJCStFS569pXLkdNgn5TKBqHt1FD9"},
		{"abandon amount liar amount expire adjust cage candy arch gather drum bullet absurd math era live bid rhythm alien crouch range attend journey unaware", "TQ8xLycC44dA9nnvME3K6X41iMjmR3J1Vz"},
	}
	g := NewTronGenerator(0)
	for _, tt := range tests {
		key, err := hdwallet.DerivePrivateKey(tt.mnemonic, hdwallet.CoinTron)
		if err != nil {
			t.Fatal(err)
		}
		if got := g.hexToBase58CheckAddress(crypto.PubkeyToAddress(key.PublicKey).Bytes()); got != tt.address {
			t.Errorf("address = %s, want %s", got, tt.address)
		}
	}
}

// 前缀 T 必然匹配，用它拿到一个钱包
func TestDoSingleMnemonic(t *testing.T) {
	for _, words := range []int{12, 24} {
		g := NewTronGenerator(words)
		wallet := g.DoSingle("T", "")
		if wallet == nil {
			t.Fatal("expected a wallet")
		}
		if n := len(strings.Fields(wallet.Mnemonic)); n != words {
			t.Errorf("mnemonic has %d words, want %d", n, words)
		}
		if wallet.Path != "m/44'/195'/0'/0/0" {
			t.Errorf("path = %s", wallet.Path)
		}

		key, err := hdwallet.DerivePrivateKey(wallet.Mnemonic, hdwallet.CoinTron)
		if err != nil {
			t.Fatal(err)
		}
		if got := hex.EncodeToString(crypto.FromECDSA(key)); got != wallet.Private {
			t.Errorf("private key %s cannot be derived from the mnemonic (%s)", wallet.Private, got)
		}
		if got := g.hexToBase58CheckAddress(crypto.PubkeyToAddress(key.PublicKey).Bytes()); got != wallet.Public {
			t.Errorf("address %s does not match the mnemonic (%s)", wallet.Public, got)
		}
	}
}

func TestDoSingleRawKey(t *testing.T) {
	wallet := NewTronGenerator(0).DoSingle("T", "")
	if wallet == nil {
		t.Fatal("expected a wallet")
	}
	if wallet.Mnemonic != "" || wallet.Path != "" {
		t.Errorf("raw key mode should not output a mnemonic: %+v", wallet)
	}
}

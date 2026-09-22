package ethereum

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"vanity-generator/cryptos/hdwallet"
)

// 前缀 0x 必然匹配，用它拿到一个钱包
func TestDoSingleMnemonic(t *testing.T) {
	for _, words := range []int{12, 24} {
		wallet := NewEthereumGenerator(words).DoSingle("0x", "")
		if wallet == nil {
			t.Fatal("expected a wallet")
		}
		if n := len(strings.Fields(wallet.Mnemonic)); n != words {
			t.Errorf("mnemonic has %d words, want %d", n, words)
		}
		if wallet.Path != "m/44'/60'/0'/0/0" {
			t.Errorf("path = %s", wallet.Path)
		}

		// 私钥和地址必须能从助记词重新派生出来
		key, err := hdwallet.DerivePrivateKey(wallet.Mnemonic, hdwallet.CoinEthereum)
		if err != nil {
			t.Fatal(err)
		}
		if got := hex.EncodeToString(crypto.FromECDSA(key)); got != wallet.Private {
			t.Errorf("private key %s cannot be derived from the mnemonic (%s)", wallet.Private, got)
		}
		if got := crypto.PubkeyToAddress(key.PublicKey).Hex(); got != wallet.Public {
			t.Errorf("address %s does not match the mnemonic (%s)", wallet.Public, got)
		}
	}
}

func TestDoSingleRawKey(t *testing.T) {
	wallet := NewEthereumGenerator(0).DoSingle("0x", "")
	if wallet == nil {
		t.Fatal("expected a wallet")
	}
	if wallet.Mnemonic != "" || wallet.Path != "" {
		t.Errorf("raw key mode should not output a mnemonic: %+v", wallet)
	}
	if len(wallet.Private) != 64 || len(wallet.Public) != 42 {
		t.Errorf("unexpected wallet: %+v", wallet)
	}
}

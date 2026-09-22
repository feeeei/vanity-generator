package hdwallet

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

// 向量由 ethers（ETH）与 tronweb（Tron）生成，路径分别为 m/44'/60'/0'/0/0 与 m/44'/195'/0'/0/0
var vectors = []struct {
	name, mnemonic     string
	ethKey, ethAddress string
	tronKey            string
}{
	{
		name:       "12 词，全零熵",
		mnemonic:   "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		ethKey:     "1ab42cc412b618bdea3a599e3c9bae199ebf030895b039e9db1e30dafb12b727",
		ethAddress: "0x9858EfFD232B4033E47d90003D41EC34EcaEda94",
		tronKey:    "b5a4cea271ff424d7c31dc12a3e43e401df7a40d7412a15750f3f0b6b5449a28",
	},
	{
		name:       "24 词，全零熵",
		mnemonic:   "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
		ethKey:     "1053fae1b3ac64f178bcc21026fd06a3f4544ec2f35338b001f02d1d8efa3d5f",
		ethAddress: "0xF278cF59F82eDcf871d630F28EcC8056f25C1cdb",
		tronKey:    "8ede55b42e315eb58ad4d98cbb3be223ae197708e015d8b635a47696f920c1a3",
	},
	{
		name:       "12 词，熵 00..0f",
		mnemonic:   "abandon amount liar amount expire adjust cage candy arch gather drum buyer",
		ethKey:     "df02cbea58239744a8a6ba328056309ae43f86fec6db45e5f782adcf38aacadf",
		ethAddress: "0x83f1caAdaBeEC2945b73087F803d404F054Cc2B7",
		tronKey:    "5e3b9f120629cd88ef264d82360e863249119aadff7d46a83ed9969524f6d5aa",
	},
	{
		name:       "24 词，熵 00..1f",
		mnemonic:   "abandon amount liar amount expire adjust cage candy arch gather drum bullet absurd math era live bid rhythm alien crouch range attend journey unaware",
		ethKey:     "e3b99b954c842dbd1e61148af4c295e15fda6d41d1633fd10fab2c6090f42669",
		ethAddress: "0xF9297b542BDb5DA50C364f9AE4Cbe1F3933bA40F",
		tronKey:    "d7f9ba29c588f6767c6df1d8c4a9b6ca737ca1ed559144a699db3f7be610414c",
	},
}

func TestDerivePrivateKey(t *testing.T) {
	for _, v := range vectors {
		t.Run(v.name, func(t *testing.T) {
			ethKey, err := DerivePrivateKey(v.mnemonic, CoinEthereum)
			if err != nil {
				t.Fatal(err)
			}
			if got := hex.EncodeToString(crypto.FromECDSA(ethKey)); got != v.ethKey {
				t.Errorf("eth key = %s, want %s", got, v.ethKey)
			}
			if got := crypto.PubkeyToAddress(ethKey.PublicKey).Hex(); got != v.ethAddress {
				t.Errorf("eth address = %s, want %s", got, v.ethAddress)
			}

			tronKey, err := DerivePrivateKey(v.mnemonic, CoinTron)
			if err != nil {
				t.Fatal(err)
			}
			if got := hex.EncodeToString(crypto.FromECDSA(tronKey)); got != v.tronKey {
				t.Errorf("tron key = %s, want %s", got, v.tronKey)
			}
		})
	}
}

func TestDerivePrivateKeyRejectsInvalidMnemonic(t *testing.T) {
	if _, err := DerivePrivateKey("abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon", CoinEthereum); err == nil {
		t.Error("mnemonic with a bad checksum should be rejected")
	}
}

func TestNewMnemonic(t *testing.T) {
	for _, words := range []int{12, 24} {
		m, err := NewMnemonic(words)
		if err != nil {
			t.Fatal(err)
		}
		if n := len(strings.Fields(m)); n != words {
			t.Errorf("NewMnemonic(%d) returned %d words", words, n)
		}
		if !bip39.IsMnemonicValid(m) {
			t.Errorf("NewMnemonic(%d) returned an invalid mnemonic: %s", words, m)
		}
	}
	for _, words := range []int{0, 11, 13, 25} {
		if _, err := NewMnemonic(words); err == nil {
			t.Errorf("NewMnemonic(%d) should fail", words)
		}
	}
}

func TestPath(t *testing.T) {
	if got := Path(CoinEthereum); got != "m/44'/60'/0'/0/0" {
		t.Errorf("eth path = %s", got)
	}
	if got := Path(CoinTron); got != "m/44'/195'/0'/0/0" {
		t.Errorf("tron path = %s", got)
	}
}

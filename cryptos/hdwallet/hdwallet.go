// Package hdwallet 提供 BIP39 助记词生成和 BIP44 路径下 secp256k1 私钥的派生，供 eth、tron 共用
package hdwallet

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/v2/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg/v2"
	"github.com/tyler-smith/go-bip39"
)

// SLIP-44 币种编号
const (
	CoinEthereum uint32 = 60
	CoinTron     uint32 = 195
)

// Path 返回钱包默认使用的 BIP44 路径 m/44'/coin'/0'/0/0，即 MetaMask、TronLink 的第一个账户
func Path(coin uint32) string {
	return fmt.Sprintf("m/44'/%d'/0'/0/0", coin)
}

// NewMnemonic 生成 words 个单词的 BIP39 英文助记词，words 必须是 12 到 24 之间 3 的倍数
func NewMnemonic(words int) (string, error) {
	if words < 12 || words > 24 || words%3 != 0 {
		return "", fmt.Errorf("unsupported mnemonic words: %d", words)
	}
	entropy, err := bip39.NewEntropy(words / 3 * 32) // 12 词 128 位，24 词 256 位
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

// NewKey 生成一个新助记词，并按 Path(coin) 派生出私钥
func NewKey(words int, coin uint32) (mnemonic string, key *ecdsa.PrivateKey, err error) {
	mnemonic, err = NewMnemonic(words)
	if err != nil {
		return "", nil, err
	}
	key, err = DerivePrivateKey(mnemonic, coin)
	if err != nil {
		return "", nil, err
	}
	return mnemonic, key, nil
}

// DerivePrivateKey 按 m/44'/coin'/0'/0/0 从助记词派生 secp256k1 私钥，不使用 BIP39 密码
func DerivePrivateKey(mnemonic string, coin uint32) (*ecdsa.PrivateKey, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	for _, index := range []uint32{hdkeychain.HardenedKeyStart + 44, hdkeychain.HardenedKeyStart + coin, hdkeychain.HardenedKeyStart, 0, 0} {
		if key, err = key.Derive(index); err != nil {
			return nil, err
		}
	}
	priv, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return priv.ToECDSA(), nil
}

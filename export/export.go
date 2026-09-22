package export

import (
	"encoding/json"
	"fmt"
	"os"
	"vanity-generator/model"
)

func Export(wallet *model.Wallet) {
	toConsole(wallet)
	toFile(wallet)
}

func toConsole(wallet *model.Wallet) {
	fmt.Println("generator wallet address:", wallet.Public)
	if wallet.Private != "" {
		fmt.Println("private_key:", wallet.Private)
	}
	if wallet.Mnemonic != "" {
		fmt.Println("mnemonic:", wallet.Mnemonic)
	}
	if wallet.Path != "" {
		fmt.Println("path:", wallet.Path)
	}
}

func toFile(wallet *model.Wallet) {
	raw, _ := json.MarshalIndent(wallet, "", "\t")
	os.WriteFile("./wallet.json", raw, 0666)
}

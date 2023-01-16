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
	fmt.Println("genertaor wallet private_key:", wallet.Private, "address:", wallet.Public)
}

func toFile(wallet *model.Wallet) {
	raw, _ := json.MarshalIndent(wallet, "", "\t")
	os.WriteFile("./wallet.json", raw, 0666)
}

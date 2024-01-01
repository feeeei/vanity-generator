package cmd

import (
	"errors"
	"strings"
	"vanity-generator/cmd/args"
	"vanity-generator/export"
	"vanity-generator/runtime"
)

func init() {
	desc := "Generate polkadot wallet address beauty address"
	addCommand("dot", desc, validatePolkadot, polkadot)
}

func validatePolkadot(args args.Args) error {
	if args.Prefix != "" && !strings.HasPrefix(args.Prefix, "1") {
		return errors.New("polkadot address should be prefixed with 1")
	}
	fix := args.Prefix + args.Suffix
	for i := range fix {
		if !isSs58(fix[i]) {
			return errors.New("address needs to follow the SS58 format")
		}
	}
	return nil
}

func polkadot(args args.Args) {
	wallet := runtime.NewExecutor("polkadot", args).Start()
	export.Export(wallet)
}

func isSs58(char byte) bool {
	switch char {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J',
		'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c',
		'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
		'w', 'x', 'y', 'z':
		return true
	default:
		return false
	}
}

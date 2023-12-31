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
	return nil
}

func polkadot(args args.Args) {
	wallet := runtime.NewExecutor("polkadot", args).Start()
	export.Export(wallet)
}

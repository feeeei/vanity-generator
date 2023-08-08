package cmd

import (
	"errors"
	"regexp"
	"strings"
	"vanity-generator/cmd/args"
	"vanity-generator/export"
	"vanity-generator/runtime"
)

func init() {
	desc := "Generate ethereum wallet address beauty address"
	addCommand("eth", desc, validateEthereum, ethereum)
}

func validateEthereum(args args.Args) error {
	if args.Prefix != "" && !strings.HasPrefix(args.Prefix, "0x") {
		return errors.New("ethereum address should be prefixed with 0x")
	}

	reg := regexp.MustCompile(`^[0-9a-fA-F]*$`)
	if len(args.Prefix) > 2 && !reg.MatchString(args.Prefix[2:]) {
		return errors.New("ethereum address format should be 0-9 or a-f or A-F")
	}
	if !reg.MatchString(args.Suffix) {
		return errors.New("ethereum address format should be 0-9 or a-f or A-F")
	}
	return nil
}

func ethereum(args args.Args) {
	wallet := runtime.NewExecutor("eth", args).Start()
	export.Export(wallet)
}

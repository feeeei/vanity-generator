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
	desc := "Generate monero wallet address beauty address"
	addCommand("xmr", desc, validateMonero, monero)
}

func validateMonero(args args.Args) error {
	if args.Prefix != "" && !strings.HasPrefix(args.Prefix, "4") {
		return errors.New("monero address should be prefixed with 4")
	}

	// 主网标准地址 4 后面的第一位只可能是 1-9、A、B
	if len(args.Prefix) > 1 && !strings.ContainsRune("123456789AB", rune(args.Prefix[1])) {
		return errors.New("monero address second place must [1-9AB]")
	}

	bodyReg := regexp.MustCompile(`^[A-HJ-NP-Za-km-z1-9]*$`)
	if !bodyReg.MatchString(args.Prefix) {
		return errors.New("monero address format should conform to the base58 specification")
	}
	if !bodyReg.MatchString(args.Suffix) {
		return errors.New("monero address format should conform to the base58 specification")
	}
	return nil
}

func monero(args args.Args) {
	wallet := runtime.NewExecutor("monero", args).Start()
	export.Export(wallet)
}

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
	desc := "Generate tron wallet address beauty address"
	addCommand("tron", desc, validateTron, tron)
}

func validateTron(args args.Args) error {
	if args.Prefix != "" && !strings.HasPrefix(args.Prefix, "T") {
		return errors.New("tron address should be prefixed with T")
	}

	// TRON 地址T后面的第一位只有9、A-HJ-NP-Z几种选项
	if len(args.Prefix) > 1 {
		if !(args.Prefix[1] == '9' || (args.Prefix[1] >= 'A' && args.Prefix[1] <= 'Z')) {
			return errors.New("tron address first place must [A-HJ-NP-Z9]")
		}
	}

	bodyReg := regexp.MustCompile(`^[A-HJ-NP-Za-km-z1-9]*$`)
	if !bodyReg.MatchString(args.Prefix) {
		return errors.New("tron address format should conform to the base58 specification")
	}
	if !bodyReg.MatchString(args.Suffix) {
		return errors.New("tron address format should conform to the base58 specification")
	}
	return nil
}

func tron(args args.Args) {
	wallet := runtime.NewExecutor("tron", args).Start()
	export.Export(wallet)
}

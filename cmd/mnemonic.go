package cmd

import (
	"errors"
	"vanity-generator/cmd/args"
)

// validateWords 校验 --words，只允许 12 或 24，0 表示使用默认值
func validateWords(args args.Args) error {
	if args.Words != 0 && args.Words != 12 && args.Words != 24 {
		return errors.New("mnemonic words must be 12 or 24")
	}
	return nil
}

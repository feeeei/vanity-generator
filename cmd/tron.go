package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
	"vanity-generator/export"
	"vanity-generator/runtime"
)

var tronCmd = &cobra.Command{
	Use:     "tron",
	Short:   "Generate tron wallet address beauty address",
	Long:    "Generate tron wallet address beauty address. Supports prefix and suffix designation",
	Run:     tron,
	PreRunE: validateTron,
}

func init() {
	rootCmd.AddCommand(tronCmd)
}

func validateTron(cmd *cobra.Command, args []string) error {
	if prefix != "" && !strings.HasPrefix(prefix, "T") {
		return errors.New("tron address should be prefixed with T")
	}

	// TRON 地址T后面的第一位只有9、A-HJ-NP-Z几种选项
	if len(prefix) > 1 {
		if !(prefix[1] == '9' || (prefix[1] >= 'A' && prefix[1] <= 'Z')) {
			return errors.New("tron address first place must [A-HJ-NP-Z9]")
		}
	}

	bodyReg := regexp.MustCompile(`^[A-HJ-NP-Za-km-z1-9]*$`)
	if !bodyReg.MatchString(prefix) {
		return errors.New("tron address format should conform to the base58 specification")
	}
	if !bodyReg.MatchString(suffix) {
		return errors.New("tron address format should conform to the base58 specification")
	}
	return nil
}

func tron(cmd *cobra.Command, args []string) {
	wallet := runtime.NewExecutor("tron", prefix, suffix, concurrency).Start()
	export.Export(wallet)
}

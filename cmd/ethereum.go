package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
	"vanity-generator/export"
	"vanity-generator/runtime"
)

var ethCmd = &cobra.Command{
	Use:     "eth",
	Short:   "Generate ethereum wallet address beauty address",
	Long:    "Generate ethereum wallet address beauty address. Supports prefix and suffix designation",
	Run:     ethereum,
	PreRunE: validateEthereum,
}

func init() {
	rootCmd.AddCommand(ethCmd)
}

func validateEthereum(cmd *cobra.Command, args []string) error {
	if prefix != "" && !strings.HasPrefix(prefix, "0x") {
		return errors.New("ethereum address should be prefixed with 0x")
	}

	reg := regexp.MustCompile(`^[0-9a-fA-F]*$`)
	if len(prefix) > 2 && !reg.MatchString(prefix[2:]) {
		return errors.New("ethereum address format should be 0-9 or a-f or A-F")
	}
	if !reg.MatchString(suffix) {
		return errors.New("ethereum address format should be 0-9 or a-f or A-F")
	}
	return nil
}

func ethereum(cmd *cobra.Command, args []string) {
	wallet := runtime.NewExecutor("eth", prefix, suffix, concurrency).Start()
	export.Export(wallet)
}

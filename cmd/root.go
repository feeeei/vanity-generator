package cmd

import (
	"github.com/spf13/cobra"
	"runtime"
)

var prefix, suffix string // 前缀、后缀匹配
var concurrency int32     // 并发控制

var rootCmd = &cobra.Command{
	Use:   "vanity-generator",
	Short: "Generate crypto wallet address beauty address",
	Long:  `Generate crypto wallet address beauty address. Supports prefix and suffix designation`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "address prefix")
	rootCmd.PersistentFlags().StringVar(&suffix, "suffix", "", "address suffix")
	rootCmd.PersistentFlags().Int32Var(&concurrency, "concurrency", int32(runtime.NumCPU()), "concurrency limit")
}

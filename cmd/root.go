package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime"
	"strings"
	"vanity-generator/cmd/args"
)

var arg args.Args

var rootCmd = &cobra.Command{
	Use:   "vanity-generator",
	Short: "Generate crypto wallet address beauty address",
	Long:  `Generate crypto wallet address beauty address. Supports prefix and suffix designation`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	initCommand()
	initConfig()
}

func initCommand() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&arg.Prefix, "prefix", "", "address prefix")
	rootCmd.PersistentFlags().StringVar(&arg.Suffix, "suffix", "", "address suffix")
	rootCmd.PersistentFlags().Int32Var(&arg.Concurrency, "concurrency", int32(runtime.NumCPU()), "concurrency limit")
}

func initConfig() {
	// file
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/vanity")
	viper.SetConfigType("json")
	viper.AllowEmptyEnv(true)
	if err := viper.ReadInConfig(); err != nil {
		if !strings.Contains(err.Error(), "Not Found") {
			log.Fatalf("read config.json error: %+v", err)
		}
	}

	// env
	viper.AutomaticEnv()

	// flag
	_ = viper.BindPFlag("prefix", rootCmd.Flag("prefix"))
	_ = viper.BindPFlag("suffix", rootCmd.Flag("suffix"))
	_ = viper.BindPFlag("concurrency", rootCmd.Flag("concurrency"))

	if len(os.Args) == 1 && viper.GetString("wallet") != "" {
		os.Args = append(os.Args, viper.GetString("wallet"))
	}
}

func addCommand(name, desc string, checkFn func(args.Args) error, runFn func(args.Args)) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   name,
		Short: desc,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			arg.Prefix = viper.GetString("prefix")
			arg.Suffix = viper.GetString("suffix")
			arg.Concurrency = viper.GetInt32("concurrency")
			return checkFn(arg)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runFn(arg)
		},
	})
}

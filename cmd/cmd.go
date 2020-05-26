package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"test.local/pkg/utils"
	"test.local/pkg/xtrace"
	"test.local/pkg/xzap"
)

var rootCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Long:  "test is a CLI library for Go that empowers applications.",
	Run: func(cmd *cobra.Command, args []string) {

	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) { // run 执行前执行，子命令若没有设置则会继承
		if isHelp() {
			return
		}

		if cfgFile == "" {
			cmd.Help()
			log.Fatal("config file nil")
		}
	},
}

var cfgFile string

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initLogger, initTrace)
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "c", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	rootCmd.AddCommand(api)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			log.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}

func initLogger() {
	if !viper.IsSet("log") {
		log.Fatalln("log config nil")
	}
	cfg := &xzap.LogConfig{
		Level:       viper.GetInt32("log.level"),
		OutputPaths: viper.GetStringSlice("log.outputPaths"),
		Dev:         viper.GetBool("log.dev"),
	}
	_, err := xzap.NewZap(utils.DefaultLoggerName, cfg)
	if err != nil {
		log.Fatalln("init logger err", err)
	}

	xzap.Async()
}

func initTrace() {
	err := xtrace.NewJaegerTracer(
		viper.GetString("name"),
		viper.GetString("log.traceLog"),
		viper.GetString("log.trace_rate"),
		"0",
	)

	if err != nil {
		log.Fatalln("init trace err", err)
	}
}

func isHelp() bool {
	for index, arg := range os.Args {
		if "help" == arg && index == 1 {
			return true
		}

		if "-h" == arg {
			return true
		}

		if "--help" == arg {
			return true
		}
	}

	return false
}

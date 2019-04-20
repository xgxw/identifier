package cmd

import (
	"fmt"
	"os"

	flog "github.com/everywan/foundation-go/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute is ..
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "identifier",
	Short: "Generate Unique Identifier",
	Long:  "支持使用 snowflake 算法为分布式系统生成全局唯一标识值",
}

func init() {
	// 读取配置
	cobra.OnInitialize(initConfig)
	// 添加配置参数
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}

func initConfig() {
	if cfgFile == "" {
		return
	}
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}
}

// ApplicationOps 程序配置文件
type ApplicationOps struct {
	Snowflake SnowflakeOps `mapstructure:"snowflake" yaml:"snowflake"`
	Logger    flog.Options `mapstructure:"logger" yaml:"logger"`
}

// Load 使用viper加载配置文件
func (opts *ApplicationOps) Load() {
	err := viper.Unmarshal(opts)
	if err != nil {
		fmt.Printf("failed to parse config file: %s", err)
	}
}

func loadApplocationOps() *ApplicationOps {
	opts := &ApplicationOps{}
	opts.Load()
	return opts
}

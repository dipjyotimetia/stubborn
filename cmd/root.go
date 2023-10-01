package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	rootCmd = &cobra.Command{
		Use:   "stubborn",
		Short: "Stubbing utility",
		Long:  "A single solution for common problems",
	}
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of stubborn",
	Long:  "All software has versions. I am stubborn",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stubborn v0.1 -- HEAD")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	rootCmd.PersistentFlags().StringP("author", "a", "Dipjyoti Metia", "author name for copyright attribution")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cmdSTUBS)
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.AddConfigPath("../configs")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func err(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

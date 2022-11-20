package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var helpz bool
var target string
var update bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tinyhelper",
	Short: "Tool for helping configure tinygo",
	Long:  `Tool for helping configure tinygo`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&helpz, "helpz", "z", false, "tinygo help")
	viper.BindPFlag("helpz", rootCmd.PersistentFlags().Lookup("helpz"))

	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "target hardware")
	viper.BindPFlag("target", rootCmd.PersistentFlags().Lookup("target"))

	rootCmd.PersistentFlags().BoolVarP(&update, "update", "u", false, "update target")
	viper.BindPFlag("update", rootCmd.PersistentFlags().Lookup("update"))
}

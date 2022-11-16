/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flashCmd represents the flash command
var flashCmd = &cobra.Command{
	Use:   "flash",
	Short: "Flash target device",
	Long: `Flash target device with passed program.
	Defaults to main function in current directory a specific path can be
	passed as an argument.

	If no target argument is passed confirms the use of previously set target
	or prompts to set new default target.

	Any tinygo flags passed will be honored.

	Exampeles:
	"tinyhelper flash"
	"tinyhelper flash ./path/to/main.go"
	"tinyhelper flash --target pico"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		runFlashCmd(args)
	},
}

func init() {
	rootCmd.AddCommand(flashCmd)
}

func runFlashCmd(args []string) {
	helpz := viper.GetBool("helpz")

	if helpz {
		out, err := exec.Command("tinygo", "flash", "--help").CombinedOutput()
		if err != nil {
			exitWithError(out)
		}

		fmt.Println(string(out))

		return
	}

	target := getTarget()

	args = append([]string{"flash", "--target", target}, args...)

	out, err := exec.Command("tinygo", args...).CombinedOutput()
	if err != nil {
		exitWithError(out)
	}

	fmt.Println(string(out))
}

func exitWithError(out []byte) {
	fmt.Println("Error:" + string(out))
	fmt.Println("Exiting...")
	os.Exit(1)
}

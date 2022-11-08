/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "configures .envrc",
	Long:  `Creates an .envrc file using data from the results of tinygo info for the passed target`,
	Run: func(cmd *cobra.Command, args []string) {
		runCmd()
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}

func runCmd() {
	fmt.Println("TinyHelper!")
	// if we are working on the tool we don't want to keep overwriting .envrc
	devMode := false
	if strings.Contains(os.Args[0], "main") {
		devMode = true
	}
	err := checkTinyGo()
	if err != nil {
		log.Fatal("Tinygo not found on $PATH. Please see https://tinygo.org/getting-started/install/ for install instructions.")
	}

	err = checkDirenv()
	if err != nil {
		log.Fatal("direnv not found on $PATH. Please see https://direnv.net/docs/installation.html for install instructions.")
	}

	target := getTarget()
	env, err := getInfo(target)
	if err != nil {
		log.Fatal(err)
	}

	i := parseInfo(env, target)

	fillTempate(i, devMode)
}

func getTarget() string {
	target := viper.GetString("target")
	if target == "" {
		target = os.Getenv("TH_TARGET")
		if target == "" {
			log.Fatal("Target required can not proceed!")
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf("No target passed. Use existing target (%s)?", target),
			Items: []string{"Yes", "No"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		if result == "No" {
			log.Fatal("Target required can not proceed, exiting!")
		}
	}

	return target
}

func getInfo(target string) (string, error) {
	out, err := exec.Command("tinygo", "info", target).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func checkTinyGo() error {
	version, err := exec.Command("tinygo", "version").Output()
	if err != nil {
		return err
	}

	fmt.Printf("TinyGo version: %s \n", version[:len(version)-1])

	return nil
}

func checkDirenv() error {
	version, err := exec.Command("direnv", "version").Output()
	if err != nil {
		return err
	}

	fmt.Printf("Direnv version: %s\n", version[:len(version)-1])

	return nil
}

type data struct {
	Goroot string
	Flags  string
	Target string
}

const (
	gorootKey = "cached GOROOT"
	tagsKey   = "build tags"
)

func parseInfo(info, target string) data {
	props := make(map[string]string)

	rows := strings.Split(info, "\n")
	for _, row := range rows {
		parts := strings.Split(row, ":")
		if len(parts) == 2 {
			props[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	flags := strings.ReplaceAll(props[tagsKey], " ", ",")

	d := data{
		Goroot: props[gorootKey],
		Flags:  fmt.Sprintf("-tags=%s", flags),
		Target: target,
	}

	return d
}

func fillTempate(info data, devMode bool) {
	tmpl, err := template.New("template").Parse("export GOROOT={{.Goroot}}\n\nexport GOFLAGS={{.Flags}}\n\nexport TH_TARGET={{.Target}}")
	if err != nil {
		log.Fatal(err)
	}
	var f *os.File
	// create the file
	if devMode == true {
		f, err = os.Create(".envrc.temp")
	} else {
		f, err = os.Create(".envrc")
	}
	if err != nil {
		log.Fatal(err)
	}
	// close the file with defer
	defer f.Close()

	tmpl.Execute(f, info)
}

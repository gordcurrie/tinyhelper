/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	yes = "Yes"
	no  = "No"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Configures .envrc",
	Long:  `Creates an .envrc file using data from the results of "tinygo info" for the passed target`,
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

	getPossibleTargets()
	target := getTarget()
	env, err := getInfo(target)
	if err != nil {
		log.Fatal(err)
	}

	i := parseInfo(env, target)

	fillTempate(i, devMode)
}

func getTarget() string {
	// "target" is a global flag
	target := viper.GetString("target")

	if target != "" {
		return target
	}

	target = os.Getenv("TH_TARGET")

	if target != "" {
		existing := promptui.Select{
			Label: fmt.Sprintf("No target passed. Use existing target (%s)?", target),
			Items: []string{yes, no},
		}
		_, result, err := existing.Run()
		if err != nil {
			log.Fatal(err)
		}

		if result == yes {
			return target
		}
	}

	choose := promptui.Select{
		Label: "Select target",
		Items: getPossibleTargets(),
	}

	_, target, err := choose.Run()
	if err != nil {
		log.Fatal(err)
	}

	return target
}

func getPossibleTargets() []string {
	out, err := exec.Command("tinygo", "targets", target).Output()
	if err != nil {
		return nil
	}

	targets := strings.Split(string(out), "\n")

	return targets
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
	gorootKey    = "cached GOROOT"
	tagsKey      = "build tags"
	commentStart = "# TinyHelper START"
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
	tmpl, err := template.New("template").Parse(commentStart + "\nexport GOROOT={{.Goroot}}\nexport GOFLAGS={{.Flags}}\nexport TH_TARGET={{.Target}}\n# TinyHelper END\n")
	if err != nil {
		log.Fatal(err)
	}

	var f *os.File
	file := ".envrc"
	if devMode == true {
		file = ".envrc.dev"
	}

	old, err := getExistingConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	f, err = os.Create(file)
	defer f.Close()

	f.Write(old)

	fmt.Printf("Writing environment config to %s...\n", file)
	tmpl.Execute(f, info)
}

func getExistingConfig(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	// regex to match everything before commentStart
	r, err := regexp.Compile(fmt.Sprintf("((.|\n)*)%s", commentStart))
	if err != nil {
		return nil, err
	}

	old := r.Find(b)

	if len(old) >= len(commentStart) { // if something matched drop commentStart from it
		old = old[:len(old)-len(commentStart)]
	} else { // no nothing matches than everything is a preexisting config
		old = b
	}

	return old, nil
}

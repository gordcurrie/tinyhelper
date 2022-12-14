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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	helpz := viper.GetBool("helpz")

	if helpz {
		out, err := exec.Command("tinygo", "help").CombinedOutput()
		if err != nil {
			exitWithError(out)
		}

		fmt.Println("tinyhelper env --helpz called")
		fmt.Println(string(out))

		return
	}

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
	if devMode {
		file = ".envrc.dev"
	}

	old, err := getExistingConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	f, err = os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = f.Write(old)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Writing environment config to %s...\n", file)
	err = tmpl.Execute(f, info)
	if err != nil {
		log.Fatal(err)
	}
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

package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func getTarget() string {
	// "target" is a global flag
	target := viper.GetString("target")

	update := viper.GetBool("update")

	if target != "" {
		return target
	}

	target = os.Getenv("TH_TARGET")

	if target != "" && !update {
		return target
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

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	bufferSize = 1024
	// should be ttyAMC.* on Linux or (tty.usbmodem.* tty.usbserial.*) on macos
	pattern = "(ttyACM.|tty\\.usbmodem.*|tty\\.usbserial.*)"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Logs serial out from /dev/ttyACM0",
	Long: `Logs serial out from /dev/ttyACM0 which should
	contain any println or fmt.Println message from you
	program.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		runLogCmd(args)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func runLogCmd(args []string) {
	conns, err := getPossibleSerialConnections()
	if err != nil {
		log.Fatalf(err.Error())
	}

	if len(conns) < 1 {
		log.Fatalf("no valid serial connections found")
	}

	var conn string

	if len(conns) > 1 {
		selectConn := promptui.Select{
			Label: "Select connection",
			Items: conns,
		}

		_, conn, err = selectConn.Run()
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		conn = conns[0]
		fmt.Printf("One valid serial connection found, logging: /dev/%s\n", conn)
	}

	f, err := os.Open("/dev/" + conn)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer f.Close()
	quit := make(chan bool, 1)

	go listenForQuit(quit)

	buf := make([]byte, bufferSize)
	loop := true
	for loop {
		n, err := f.Read(buf)
		if err == io.EOF {
			log.Println("Connection lost, exiting...")
			break
		}

		if err != nil {
			log.Printf("%#v\n", err)
			continue
		}

		out := clean(buf[:n])

		if n > 0 && out != "" {
			fmt.Printf("Output: %v\n", out)
		}

		select {
		case exit := <-quit:
			loop = !exit
		default:
			continue
		}
	}
}

func getPossibleSerialConnections() ([]string, error) {
	out, err := exec.Command("ls", "/dev/").CombinedOutput()
	if err != nil {
		return nil, err
	}
	files := strings.Split(string(out), "\n")

	var results []string
	for _, file := range files {
		match, err := regexp.MatchString(pattern, file)
		if err != nil {
			continue
		}
		if match {
			results = append(results, file)
		}
	}

	return results, nil
}

func listenForQuit(quit chan bool) {
	reader := bufio.NewReader(os.Stdin)

	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			return
		}

		if strings.Contains(string(in), ":q") {
			fmt.Println("Quit called, exiting...")
			quit <- true
		}
	}
}

func clean(buf []byte) string {
	return strings.Replace(string(buf), "\n", "", -1)
}

package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

var underline = color.New(color.Underline).SprintFunc()
var bold = color.New(color.Bold).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func checkBinaryExists(command string) (bool, string) {
	out := outputCommand(fmt.Sprintf("which %s", command))
	return out != "", out
}

func outputCommand(command string) string {
	out, _ := exec.Command("sh", "-c", command).Output()

	return strings.TrimSuffix(string(out), "\n")
}

func runCommand(command string) {
	exec.Command("sh", "-c", command).Run()

	return
}

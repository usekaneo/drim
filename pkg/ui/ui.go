package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

func Info(message string) {
	fmt.Printf("%s%s%s\n", ColorBlue, message, ColorReset)
}

func Success(message string) {
	fmt.Printf("%s%s%s\n", ColorGreen, message, ColorReset)
}

func Warning(message string) {
	fmt.Printf("%s%s%s\n", ColorYellow, message, ColorReset)
}

func Error(message string) {
	fmt.Printf("%s%s%s\n", ColorRed, message, ColorReset)
}

func Prompt(message string) (string, error) {
	fmt.Printf("%s%s%s ", ColorCyan, message, ColorReset)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func Confirm(message string) (bool, error) {
	response, err := Prompt(fmt.Sprintf("%s (y/N)", message))
	if err != nil {
		return false, err
	}
	response = strings.ToLower(response)
	return response == "y" || response == "yes", nil
}

func Choose(message string, options []string) (string, error) {
	fmt.Printf("%s%s%s\n", ColorCyan, message, ColorReset)
	for i, option := range options {
		fmt.Printf("%s[%d] %s%s\n", ColorCyan, i+1, option, ColorReset)
	}
	choice, err := Prompt("Enter the number of your choice:")
	if err != nil {
		return "", err
	}
	index, err := strconv.Atoi(choice)
	if err != nil || index < 1 || index > len(options) {
		return "", fmt.Errorf("invalid choice")
	}
	return options[index-1], nil
}

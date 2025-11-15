package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorReset   = "\033[0m"
	clearDisplay = "\033[H\033[2J"
)

// Очистка экрана

func clearScreen() {
	fmt.Print(clearDisplay, colorReset)
}

// Вывод сообщения об успехе

func showSuccess(message string) {
	fmt.Printf("%s✓ Success: %s%s\n", colorGreen, message, colorReset)
}

// Вывод сообщения об ошибке

func showError(message string) {
	fmt.Printf("%s✗ Error: %s%s\n", colorRed, message, colorReset)
}

// Вывод информационного сообщения

func showInfo(message string) {
	fmt.Printf("%s→ Info: %s%s\n", colorYellow, message, colorReset)
}

// Ожидание нажатия Enter

func waitForEnter() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press Enter to continue...")
	reader.ReadString('\n')
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// Алгоритм работы
// 1. Показать приглашение к вводу
// 2. Прочитать строку до символа новой строки
// 3. Удалить пробельные символы в начале и конце
// 4. Вернуть очищенную от боковых пробелов строку

func ReadUserInput(prompt string) string {

	// 1
	fmt.Print(prompt)

	// 2
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения ввода:", err)
		return ""
	}

	// 3
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		fmt.Println("Ввод не может быть пустым! Попробуйте снова.")
		return ""
	}
	// 4
	return trimmed

}

// Алгоритм работы
//
// 1. Отключить эхо-вывод в терминале
// 2. Прочитать пароль
// 3. Восстановить нормальный режим терминала
// 4. Добавить перевод строки (так как его нет при скрытом вводе)
// 5. Вернуть введённый пароль

func readPassword() (string, error) {

	fd := int(os.Stdin.Fd())

	// 1
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	// 3
	defer term.Restore(fd, oldState)

	// 2
	passByte, err := term.ReadPassword(fd)
	if err != nil {
		return "", err
	}
	pass := string(passByte)

	// 4
	fmt.Println()

	// 5
	return pass, nil
}

// Алгоритм работы
//
// 1. Очистить экран
// 2. Показать заголовок
// 3. Вывести список всех доступных команд
// 4. Добавить разделители для лучшей читаемости

func ShowMainMenu() {
	clearScreen()

	// Повторяем символ равенства n кол-во раз
	sepLine := strings.Repeat("=", 42)
	fmt.Println(sepLine)

	title := "Password Manager"
	// Центрируем заголовок длиной 42 символа
	padding := (42 - len(title)) / 2
	fmt.Printf("%s%s%s\n", strings.Repeat(" ", padding), title, strings.Repeat(" ", padding))

	fmt.Println(sepLine)

	commands := []string{
		"1. Generate new password",
		"2. Add new password",
		"3. Get password",
		"4. List all passwords",
		"5. Update password",
		"6. Delete password",
		"7. List categories",
		"8. Show password statistics",
		"9. Find duplicate passwords",
		"0. Exit",
	}

	for _, cmd := range commands {
		fmt.Println(cmd)
	}

	// 4. Разделитель
	fmt.Println(sepLine + "\n")
}

// Алгоритм работы
//
// 1. Вывести заголовки столбцов
// 2. Создать разделительную линию
// 3. Вывести данные каждого пароля в табличном формате
// 4. Скрыть значения паролей для безопасности

func PrintPasswordList(passwords []Password) {

	fmt.Printf("%-20s %-15s %-20s %-20s\n", "Name", "Category", "Created", "Last Modified")
	fmt.Println(strings.Repeat("-", 80))

	for _, v := range passwords {
		fmt.Printf("%-20s %-15s %-20s %-20s\n", v.Name, v.Category, v.CreatedAt.Format("2006-01-02"), v.LastModified.Format("2006-01-02"))
	}

}

// Алгоритм работы
//
// 1. Показать все поля пароля
// 2. Отформатировать даты в читаемом формате
// 3. Структурировать вывод для лучшей читаемости

func ShowPasswordDetails(password Password) {
	fmt.Printf("Service: %s\n", password.Name)
	fmt.Printf("Category: %s\n", password.Category)
	fmt.Printf("Password: %s\n", password.Value)
	fmt.Printf("Created: %s\n", password.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Last Modified: %s\n", password.LastModified.Format("2006-01-02 15:04:05"))
}

// В обработчиках используется готовая функция passInput, чтобы не повторять один и тот же код
func passInput(pm *PasswordManager) (string, error) {
	fmt.Print("Enter password (or press Enter to generate): ")
	passIn, err := readPassword()
	if err != nil {
		return "", err
	}

	if passIn == "" {
		clearScreen()
		input := ReadUserInput("Enter password length (min 8): ")
		length, err := strconv.Atoi(input)
		if err != nil {
			showError("Invalid number")
			return "", err
		}

		pass, err := pm.GeneratePassword(length)
		if err != nil {
			return "", err
		}

		passIn = pass

		showInfo(fmt.Sprintf("Generated password: %s", passIn))
	} else {
		clearScreen()
		err := pm.CheckPasswordStrength(passIn)
		if err != nil {
			return "", ErrPassWeak
		}
		showInfo(fmt.Sprintf("Generated password: %s", passIn))
	}

	return passIn, nil
}

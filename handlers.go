package main

import (
	"fmt"
	"strconv"
	"time"
)

// Алгоритм работы
//
// 1. Запросить длину пароля
// 2. Сгенерировать пароль
// 3. Показать результат
// 4. Обработать возможные ошибки

func HandlePasswordGeneration(pm *PasswordManager) error {
	clearScreen()

	input, err := ReadUserInput("Enter password length (min 8): ")
	if err != nil {
		return err
	}

	length, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	pass, err := pm.GeneratePassword(length)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	showSuccess("Password generated successfully")
	fmt.Println("Generated password:", pass)

	fmt.Println()

	waitForEnter()

	return nil
}

// Алгоритм работы
//
// 1. Запросить имя сервиса
// 2. Предложить ввести пароль или сгенерировать новый
// 3. Запросить категорию
// 4. Сохранить пароль
// 5. Показать результат операции

func HandlePasswordAdd(pm *PasswordManager) error {
	clearScreen()
	nameInput, err := ReadUserInput("Enter service name: ")
	if err != nil {
		return err
	}

	input, err := passInput(pm)
	if err != nil {
		return err
	}

	clearScreen()
	catInput, err := ReadUserInput("Enter category: ")
	if err != nil {
		return err
	}

	if err = pm.SavePassword(nameInput, input, catInput); err != nil {
		return err
	}

	showSuccess("Password saved successfully\n")

	waitForEnter()

	return nil
}

// Алгоритм работы
//
// 1. Запросить имя сервиса
// 2. Найти пароль
// 3. Показать детальную информацию
// 4. Обработать случай отсутствия пароля

func HandlePasswordSearch(pm *PasswordManager) error {
	clearScreen()
	nameInput, err := ReadUserInput("Enter service name: ")
	if err != nil {
		return err
	}

	pass, err := pm.GetPassword(nameInput)
	if err != nil {
		return err
	}

	fmt.Println("Password Details:")
	fmt.Println("Service:", pass.Name)
	fmt.Println("Category:", pass.Category)
	fmt.Println("Password:", pass.Value)
	fmt.Println("Created:", pass.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("Last Modified:", pass.LastModified.Format("2006-01-02 15:04:05"))

	fmt.Println()

	waitForEnter()

	return nil
}

// Алгоритм работы
//
// 1. Запросить имя сервиса
// 2. Запросить новый пароль
// 3. Обновить запись
// 4. Показать результат обновления

func HandlePasswordUpdate(pm *PasswordManager) error {
	clearScreen()
	nameInput, err := ReadUserInput("Enter service name: ")
	if err != nil {
		return err
	}

	newValue, err := passInput(pm)
	if err != nil {
		return err
	}

	if err = pm.UpdatePassword(nameInput, newValue); err != nil {
		return err
	}

	showSuccess("Password updated successfully!\n")

	waitForEnter()

	return nil
}

//Алгоритм работы функции:
//
// 1. Очистить экран и показать сообщение о процессе сохранения
// 2. Попытаться сохранить данные в файл
// 3. Показать результат операции (успех или ошибка)
// 4. Вывести прощальное сообщение при успешном сохранении

func HandleExitAndSave(pm *PasswordManager) error {
	clearScreen()
	fmt.Println("Saving changes...")
	if err := pm.SaveToFile(); err != nil {
		return err
	}

	showSuccess("Changes saved successfully!")
	showSuccess("Goodbye!")

	waitForEnter()

	return nil
}

// Алгоритм работы
//
// 1. Пройти по всем элементам []Password и вывести значения
// 2. Обработать ошибки

func HandlePasswordsList(pm *PasswordManager) error {
	clearScreen()

	passwords := pm.ListPasswords()
	fmt.Printf("Total passwords: %d\n\n", len(passwords))
	for _, p := range passwords {
		fmt.Printf("Service: %-25s Category: %-15s CreatedAt: %s\n", p.Name, p.Category, p.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	fmt.Println()
	waitForEnter()

	return nil
}

// Алгоритм работы
//
// 1. Ввести имя сервиса
// 2. Удалить пароль
// 3. Обработать ошибки

func HandlePasswordDelete(pm *PasswordManager) error {
	clearScreen()

	nameInput, err := ReadUserInput("Enter service name: ")
	if err != nil {
		return err
	}

	if err := pm.DeletePassword(nameInput); err != nil {
		return err
	}

	showSuccess("Password deleted successfully\n")

	waitForEnter()

	return nil
}

func HandlePasswordListCategories(pm *PasswordManager) error {
	clearScreen()

	categories := pm.ListCategories()
	fmt.Printf("Total categories: %d\n\n", len(categories))
	fmt.Println("List of categories:")
	for _, category := range categories {
		count := len(pm.GetPasswordsByCategory(category))
		fmt.Printf("- %s (%d passwords)\n", category, count)
	}

	fmt.Println()
	waitForEnter()

	return nil
}

func HandlePasswordStats(pm *PasswordManager) error {
	clearScreen()

	stats := pm.GetPasswordStats()

	fmt.Printf("Total statistics:\n")
	fmt.Printf("⚡ Total passwords:   %d\n", stats["total_passwords"])

	fmt.Printf("\n📂 Distribution by categories:\n")
	if categories, ok := stats["categories"].(map[string]int); ok {
		for category, count := range categories {
			fmt.Printf("   • %-15s: %d\n", category, count)
		}
	}

	if oldestDate, ok := stats["oldest_password_date"].(time.Time); ok {
		fmt.Printf("\n🕒 Time characteristics:\n")
		fmt.Printf("   • Oldest: %s\n", oldestDate.Format("2006-01-02"))
		if newestDate, ok := stats["newest_password_date"].(time.Time); ok {
			fmt.Printf("   • Newest: %s\n", newestDate.Format("2006-01-02"))
		}
	}

	fmt.Println()
	waitForEnter()

	return nil
}

func HandlePasswordDuplicate(pm *PasswordManager) error {
	clearScreen()

	duplicates := pm.FindDuplicatePasswords()

	if len(duplicates) == 0 {
		fmt.Println("Duplicates not found")
	} else {
		fmt.Printf("\nFound duplicates:\n")
		for password, services := range duplicates {
			fmt.Printf("\nPassword '%s' is used in the following services:\n", password)
			for _, service := range services {
				fmt.Printf("- %s\n", service)
			}
		}
	}

	fmt.Println()
	waitForEnter()

	return nil
}

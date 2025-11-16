package main

import (
	"fmt"
	"os"
)

// 1.  Реализуем структуру Password (pass.go)
// 2.  Структура PasswordManager (pass.go)
// 3.  Функция GeneratePassword (pass.go)
// 4.  Функция SavePassword (pass.go)
// 5.  Функция GetPassword (pass.go)
// 6.  Функция ListPasswords (pass.go)
// 7.  Функция SetMasterPassword (pass.go)
// 8.  Функция SaveToFile (file.go)
// 9.  Функция LoadFromFile (file.go)
// 10. Функция CheckPasswordStrength (pass.go)
// 11. Функция GetPasswordByCategory (category.go)
// 12. Функция FindDuplicatePasswords (pass.go)
// 13. Функция UpdatePassword (pass.go)
// 14. Функция DeletePassword (pass.go)
// 15. Функция ListCategories (category.go)
// 16. Функция GetPasswordStats (pass.go)
// 17. Базовые UI функции (ui.go)
// 18. Функции ввода (app.go)
// 19. Функции отображения (app.go)
// 20. Функции обработки команд (handlers.go)
// 21. Функция HandleExitAndSave (handlers.go)
// 22. Main()

func main() {

	clearScreen()
	pm := NewPasswordManager("ne_password.dat")

	fmt.Println("=== Password Manager Initialization ===")
	fmt.Print("Enter master password: ")
	masterPassword, err := readPassword()
	if err != nil {
		showError(fmt.Sprintf("Error reading master password: %v", err))
		waitForEnter()
		return
	}
	clearScreen()
	if err := pm.SetMasterPassword(masterPassword); err != nil {
		showError(fmt.Sprintf("Error setting master password: %v", err))
		waitForEnter()
		return
	}

	// Пытаемся загрузить существующие данные
	if err := pm.LoadFromFile(); err != nil && !os.IsNotExist(err) {
		showError(fmt.Sprintf("Error loading data: %v", err))
		waitForEnter()
		return
	}

	showSuccess("Password manager initialized successfully")
	waitForEnter()

	for {
		ShowMainMenu()
		var err error

		choice, err := ReadUserInput("Enter your choice: ")

		switch choice {
		case "1":
			err = HandlePasswordGeneration(pm)
		case "2":
			err = HandlePasswordAdd(pm)
		case "3":
			err = HandlePasswordSearch(pm)
		case "4":
			err = HandlePasswordsList(pm)
		case "5":
			err = HandlePasswordUpdate(pm)
		case "6":
			err = HandlePasswordDelete(pm)
		case "7":
			err = HandlePasswordListCategories(pm)
		case "8":
			err = HandlePasswordStats(pm)
		case "9":
			err = HandlePasswordDuplicate(pm)
		case "0":
			clearScreen()
			fmt.Println("=== Saving and Exiting ===")
			err = HandleExitAndSave(pm)
			if err != nil {
				showError(fmt.Sprintf("Error during exit: %v", err))
				waitForEnter()
				return
			}
			return
		default:
			showError("Invalid choice. Please try again")
			waitForEnter()
		}

		if err != nil {
			showError(fmt.Sprintf("Operation failed: %v", err))
			waitForEnter()
		}
	}
}

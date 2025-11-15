package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
	"os"
)

// Алгоритм работы функции:
//
// 1. Проверить, что менеджер инициализирован
// 2. Сериализовать map паролей в JSON
// 3. Создать новый блок шифрования AES
// 4. Сгененрировать случайный вектор инициализации
// 5. Создать шифровальщик и зашифровать данные
// 6. Сохранить IV и зашифрованные данные в файл

func (pm *PasswordManager) SaveToFile() error {
	pm.mu.RLock()

	// 1
	if pm.isInitialized != true {
		return ErrPassManagerNotInit
	}

	// 2
	data, err := json.Marshal(pm.passwords)
	if err != nil {
		return err
	}

	pm.mu.RUnlock()

	// 3
	block, err := aes.NewCipher(pm.masterKey)
	if err != nil {
		return err
	}

	// 4
	iv := make([]byte, aes.BlockSize) // BlockSize = 16
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	// 5
	stream := cipher.NewCFBEncrypter(block, iv)
	encryptedData := make([]byte, len(data))
	stream.XORKeyStream(encryptedData, data)

	// 6
	file, err := os.Create(pm.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Сначала записываем IV
	if _, err := file.Write(iv); err != nil {
		return err
	}

	// Затем шифрованные данные
	if _, err := file.Write(encryptedData); err != nil {
		return err
	}

	return nil
}

// Алгоритм работы функции:
//
// 1. Проверить, что менеджер инициализирован
// 2. Открыть файл и прочитать его содержимое
// 3. Прочитать вектор инициализации (первые 16 байт)
// 4. Прочитать зашифрованные данные
// 5. Создать расшифровщик и расшифровать данные
// 6. Преобразовать расшифрованные данные обратно в структуры

func (pm *PasswordManager) LoadFromFile() error {

	// 2
	file, err := os.Open(pm.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 3
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(file, iv); err != nil {
		return err
	}

	// 4
	encryptedData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// 5
	block, err := aes.NewCipher(pm.masterKey)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	decryptedData := make([]byte, len(encryptedData))
	stream.XORKeyStream(decryptedData, encryptedData)

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 1
	if pm.isInitialized != true {
		return ErrPassManagerNotInit
	}

	// 6
	return json.Unmarshal(decryptedData, &pm.passwords)
}

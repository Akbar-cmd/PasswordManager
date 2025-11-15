package main

import (
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Password struct {
	// Название сервиса или сайта
	Name string `json:"name"`
	// Значение пароля
	Value string `json:"value"`
	// Категория для группировки("social", "work", "finance")
	Category string `json:"category"`
	// Дата создания записи
	CreatedAt time.Time `json:"createdAt"`
	// Дата последнего изменения
	LastModified time.Time `json:"lastModified"`
}

func NewPassword(name, value, category string) *Password {
	now := time.Now()
	return &Password{
		Name:         name,
		Value:        value,
		Category:     category,
		CreatedAt:    now,
		LastModified: now,
	}
}

type PasswordManager struct {
	// Хранилище паролей, где ключ - название сервиса
	passwords map[string]Password `json:"passwords"`
	// Главный ключ шифрования, используется для защиты всех паролей
	masterKey []byte `json:"-"`
	// Путь к файлу для хранения зашифрованных данных
	filePath string `json:"-"`
	// Флаг, показывающий установлен ли мастер-пароль
	isInitialized bool `json:"-"`
	// (ОТ себя) добавил mutex
	mu sync.RWMutex
}

func NewPasswordManager(filePath string) *PasswordManager {
	return &PasswordManager{
		passwords:     make(map[string]Password),
		masterKey:     make([]byte, 0),
		filePath:      filePath,
		isInitialized: false,
	}
}

// Алгоритм работы функции:
//
//	1.Проверить, что длина пароля не меньше 8 символов
//	2.Определить набор допустимых символов (буквы, цифры, специальные символы)
//	3.Использовать crypto/rand для генерации случайных байтов
//	4.Преобразовать случайные байты в символы из допустимого набора
//	5.Вернуть сгенерированный пароль или ошибку

func (pm *PasswordManager) GeneratePassword(length int) (string, error) {
	// 1
	if length < 8 {
		return "", fmt.Errorf("password length must be at least 8, got %d", length)
	}

	// 2
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*-_+=."

	// 3
	cont := make([]byte, length)
	_, err := rand.Read(cont)
	if err != nil {
		return "", err
	}

	for i := 0; i <= length-1; i++ {
		index := cont[i] % byte(len(charset))
		cont[i] = charset[index]
	}

	return string(cont), nil

}

// Алгоритм работы функции:
//
// 1.Проверить, что менеджер паролей инициализирован (isInitialized == true)
// 2.Убедиться, что пароль с таким именем ещё не существует в хранилище
// 3.Создать новую запись пароля с помощью функции NewPassword
// 4.Сохранить пароль в map хранилища
// 5.Вернуть ошибку, если что-то пошло не так

func (pm *PasswordManager) SavePassword(name, value, category string) error {
	// на всякий
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 1
	if pm.isInitialized != true {
		return ErrPassManagerNotInit
	}

	// 2
	if _, ok := pm.passwords[name]; ok {
		return ErrPassExists
	}

	// 3
	pass := NewPassword(name, value, category)

	// 4
	pm.passwords[name] = *pass

	// 5
	return nil
}

//Алгоритм работы функции:
//
//1.Проверить, что менеджер паролей инициализирован
//2.Найти пароль в хранилище по имени
//3.Вернуть найденный пароль или ошибку, если пароль не найден

func (pm *PasswordManager) GetPassword(name string) (Password, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// 1
	if pm.isInitialized != true {
		return Password{}, ErrPassManagerNotInit
	}

	// 2
	if _, ok := pm.passwords[name]; ok {
		return pm.passwords[name], nil
	}
	// 3
	return Password{}, ErrPassNotFound
}

//Алгоритм работы функции:
//
//1. Создать новый слайс для хранения паролей
//2. Пройти по всем паролям в map хранилища
//3. Добавить каждый пароль в результирующий слайс
//4. Вернуть слайс со всеми паролями

func (pm *PasswordManager) ListPasswords() []Password {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// 1
	// создаем слайс длиной 0 и емкостью длины кол-ва паролей
	res := make([]Password, 0, len(pm.passwords))

	// 2
	for _, v := range pm.passwords {
		res = append(res, v)
	}

	return res
}

//Алгоритм работы функции:
//
//Проверить длину мастер-пароля (минимум 8 символов)
//Создать байтовый слайс размером 32 байта
//Скопировать байты мастер-пароля в этот слайс
//Сохранить полученный слайс в поле masterKey
//Установить флаг isInitialized в true

func (pm *PasswordManager) SetMasterPassword(masterPassword string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 1
	if len(masterPassword) < 8 {
		return ErrPassWeak
	}

	// 2
	contByte := make([]byte, 32)

	// 3
	copy(contByte, masterPassword)

	// 4
	pm.masterKey = contByte

	// 5
	pm.isInitialized = true

	return nil
}

// Алгоритм работы функции:
//
// 1. Проверить минимальную длину пароля (не менее 8 символов)
// 2. Проверить наличие символов разных категорий:
// 3. Убедиться, что пароль содержит символы всех категорий
// 4. Вернуть ошибку, если какое-либо требование не выполнено

func (pm *PasswordManager) CheckPasswordStrength(password string) error {
	// 1
	if len(password) < 8 {
		return ErrPassWeak
	}

	// 2
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*", char):
			hasSpecial = true
		}
	}

	// 3
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrPassWeak
	}

	// 4
	return nil
}

//Алгоритм работы функции:
//
// 1. Создать карту для хранения дубликатов, где:
//		ключ - значение пароля
//		значение - список имён сервисов, использующих этот пароль
// 2. Перебрать все пары паролей в хранилище
// 3. Если значения паролей совпадают, добавить их в карту дубликатов
// 4. Вернуть найденные дубликаты

func (pm *PasswordManager) FindDuplicatePasswords() map[string][]string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// 1
	// результирующая map
	res := make(map[string][]string)
	// для группировки сервисов по паролям
	dblPass := make(map[string][]string)

	// 2
	for k, v := range pm.passwords {
		// берем значение пароля как ключ и добавляем в промежуточную map
		dblPass[v.Value] = append(dblPass[v.Value], k)
	}

	// 3
	for k, v := range dblPass {
		// если дубликат
		if len(v) > 1 {
			// добавляем в результирующую таблицу
			res[k] = v
		}
	}

	// 4
	return res
}

// Алгоритм работы функции:
//
// 1. Проверить, что менеджер инициализирован
// 2. Найти пароль в хранилище по имени
// 3. Проверить надёжность нового пароля через CheckPasswordStrength
// 4. Обновить значение пароля
// 5. Обновить время последнего изменения
// 6. Сохранить обновлённую запись в хранилище

func (pm *PasswordManager) UpdatePassword(name, newValue string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 1
	if pm.isInitialized != true {
		return ErrPassManagerNotInit
	}

	// 2
	if _, ok := pm.passwords[name]; !ok {
		return ErrPassNotFound
	}

	// 3
	if err := pm.CheckPasswordStrength(newValue); err != nil {
		return ErrPassWeak
	}

	// 4
	// получаем копию структуры Password
	p := pm.passwords[name]
	p.Value = newValue

	// 5
	p.LastModified = time.Now()

	// Записываем изменения в структуру Password
	pm.passwords[name] = p

	// 6
	return nil
}

//Алгоритм работы функции:
//
//Проверить, что менеджер инициализирован
//Проверить существование пароля в хранилище
//Удалить запись из map хранилища
//Вернуть ошибку, если что-то пошло не так

func (pm *PasswordManager) DeletePassword(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 1
	if pm.isInitialized != true {
		return ErrPassManagerNotInit
	}

	// 2
	if _, ok := pm.passwords[name]; !ok {
		return ErrPassNotFound
	}

	// 3
	delete(pm.passwords, name)

	// 4
	return nil
}

// Алгоритм работы функции:
//
// 1. Создать карту для статистики
// 2. Собрать базовую информацию:
// 		Общее количество паролей
// 		Подсчёт паролей по категориям
// 3. Найти временные характеристики:
// 		Дата самого старого пароля
// 		Дата самого нового пароля
// 4. Сохранить все метрики в карту статистики

func (pm *PasswordManager) GetPasswordStats() map[string]interface{} {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 1
	stats := make(map[string]interface{})
	countCat := make(map[string]int)
	oldTime := time.Time{}
	newTime := time.Time{}
	// 2
	countPass := len(pm.passwords)

	for _, v := range pm.passwords {
		countCat[v.Category]++
	}

	// 3
	for _, v := range pm.passwords {
		// isZero() - проверяет нулевое ли значение
		if oldTime.IsZero() || v.CreatedAt.Before(oldTime) {
			oldTime = v.CreatedAt
		}

		if newTime.IsZero() || v.CreatedAt.After(newTime) {
			newTime = v.CreatedAt
		}
	}

	// 4
	stats["total_passwords"] = countPass
	stats["categories"] = countCat
	stats["oldest_password_date"] = oldTime
	stats["newest_password_date"] = newTime

	return stats
}

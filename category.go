package main

// Алгоритм работы функции:
//
// 1. Создать пустой слайс для хранения результатов
// 2. Пройти по всем паролям в хранилище
// 3. Сравнить категории и если совпадает добавить в результат
// 4. Вернуть список найденных паролей

func (pm *PasswordManager) GetPasswordsByCategory(category string) []Password {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// 1
	res := make([]Password, 0, len(pm.passwords))

	// 2
	for _, v := range pm.passwords {
		// 3
		if v.Category == category {
			res = append(res, v)
		}

	}
	// 4
	return res
}

//Алгоритм работы функции:
//
// 1. Создать карту для отслеживания уникальных категорий:
// 2. Пройти по всем паролям и добавить их категории в карту
// 3. Преобразовать ключи карты в слайс строк
// 4. Вернуть список уникальных категорий

func (pm *PasswordManager) ListCategories() []string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 1
	// Создали мапу для отслеживания категорий, если true - значит категория найдена
	categories := make(map[string]bool)
	// Создали результирующий слайс для возврата данных
	res := make([]string, 0, len(pm.passwords))

	// 2
	for _, v := range pm.passwords {
		// Получаем категорию и добавляем ее в map как ключ со значением true
		categories[v.Category] = true
	}

	// 3
	for k := range categories {
		// ключ (категорию) добавляем в слайс
		res = append(res, k)
	}

	// 4
	return res
}

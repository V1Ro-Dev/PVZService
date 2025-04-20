package utils

import (
	"os"
)

// GetEnv считывает значение переменной окружения или возвращает значение по умолчанию, если переменная не установлена
func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

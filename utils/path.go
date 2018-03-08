package utils

import (
	"os"
)

// Getwd get project dir
func Getwd() (string, error) {
	return os.Getwd()
}

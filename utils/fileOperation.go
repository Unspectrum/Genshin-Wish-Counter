package utils

import (
	"io"
	"os"
	"strings"
)

// Internal function to open file
func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Internal function to read from file with certain delimiter
// This function return empty string and error object if any error happen
func readFileToStringArray(file *os.File, delim string) ([]string, error) {
	dataBytes, err := io.ReadAll(file)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(dataBytes), delim), nil
}

// Function to open file from filepath and split by delim.
// This function return empty string and error object if any error happen
func OpenFileToStringArray(filepath string, delim string) ([]string, error) {
	file, err := openFile(filepath)
	if err != nil {
		return []string{}, err
	}
	lines, err := readFileToStringArray(file, delim)
	return lines, nil
}

func OpenReadFileToString(str string) (string, error) {
	b, err := os.ReadFile(str)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

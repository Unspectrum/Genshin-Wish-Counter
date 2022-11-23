package utils

import (
	"io"
	"os"
	"strings"
)

const (
	dataFileLocation = "webCaches/Cache/Cache_Data/data_2"
	warmUpStr        = "Warmup file "
	streamAssetsStr  = "StreamingAssets"
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

func ParseInstallLocation(str string) string {
	str = strings.ReplaceAll(str, warmUpStr, "")
	str = strings.Split(str, "\\")[0]
	return strings.ReplaceAll(str, streamAssetsStr, "")
}

func GetDataFileLocation(str []string) string {
	var installLocation string
	for _, line := range str {
		if strings.Contains(line, warmUpStr) {
			installLocation = ParseInstallLocation(line)
			break
		}
	}
	installLocation += dataFileLocation

	return installLocation
}

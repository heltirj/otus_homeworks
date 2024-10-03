package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrInvalidFileName     = errors.New("invalid file name")
	ErrFailedToGetFileSize = errors.New("failed to get file size")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirItems, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for i := range dirItems {
		if dirItems[i].IsDir() {
			continue
		}

		val, err := parseFile(dir, dirItems[i].Name())
		if err != nil {
			return nil, err
		}

		env[dirItems[i].Name()] = val
	}

	return env, nil
}

func parseFile(dirPath string, fileName string) (EnvValue, error) {
	val := EnvValue{}
	if strings.Contains(fileName, "=") {
		return EnvValue{}, fmt.Errorf("%w: %s", ErrInvalidFileName, fileName)
	}

	f, err := os.Open(dirPath + "/" + fileName)
	if err != nil {
		return EnvValue{}, err
	}
	defer f.Close()

	stats, err := f.Stat()
	if err != nil {
		return EnvValue{}, ErrFailedToGetFileSize
	}

	if stats.Size() == 0 {
		return EnvValue{NeedRemove: true}, nil
	}

	reader := bufio.NewReader(f)
	content, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return EnvValue{}, err
	}

	content = bytes.Replace(content, []byte{0}, []byte{'\n'}, -1)
	val.Value = strings.TrimRight(string(content), "\n\t ")

	return val, nil
}

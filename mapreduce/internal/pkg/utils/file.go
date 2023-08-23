package utils

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type Line []byte

func ReadFile(filename string) ([]Line, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "read file %s error", filename)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []Line
	for scanner.Scan() {
		lines = append(lines, scanner.Bytes())
	}

	return lines, scanner.Err()
}

func SplitLinesToFile(lines []Line, n int, prefix string) ([]string, error) {
	if len(lines) == 0 {
		return nil, errors.New("lines is empty")
	}

	fileNum := len(lines) / n
	var res []string
	for i := 0; i*fileNum < len(lines); i++ {
		fileName := fmt.Sprintf("/tmp/%s-%d.cache", prefix, i)
		size := 0
		start := i * fileNum
		for j := start; j < start+fileNum && j < len(lines); j++ {
			size += len(lines[j])
		}

		contents := make([]byte, 0, size+fileNum)
		for j := start; j < start+fileNum && j < len(lines); j++ {
			contents = append(contents, lines[j]...)
			contents = append(contents, '\n')
		}

		res = append(res, fileName)
		if err := os.WriteFile(fileName, contents, 0644); err != nil {
			return nil, errors.Wrapf(err, "write file %s error", fileName)
		}
	}

	return res, nil
}

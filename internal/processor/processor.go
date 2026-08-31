package processor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Processor struct {
	currentYear  int
	currentMonth time.Month
	currentDay   int
}

type MoolaFile [][]string

func ReadMoolaFile(path string) (MoolaFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result MoolaFile
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())

		if len(words) > 0 {
			result = append(result, words)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func ProcessMoolaFile(path string, isMasterFile bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var mf MoolaFile
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())

		if len(words) > 0 {
			mf = append(mf, words)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, line := range mf {
		fmt.Println(line[0])
	}

	return nil
}

package proloc

import (
	"bufio"
	"os"
)

func calculateLinesCount(file string) (uint64, error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	lineCount := 0
	for {
		_, err := reader.ReadString('n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}
		lineCount += 1
	}

	return uint64(lineCount), nil
}

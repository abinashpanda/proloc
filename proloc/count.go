package proloc

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

func calculateLinesCount(file string) uint64 {
	f, err := os.Open(file)
	if err != nil {
		return 0
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

	return uint64(lineCount)
}

var (
	THOUSAND uint64 = uint64(math.Pow(10, 3))
)

func formatNumber(num uint64) string {
	result := []string{}
	val := num
	for val > 0 {
		remainder := val % THOUSAND
		result = append([]string{fmt.Sprintf("%d", remainder)}, result...)
		val = val / THOUSAND
	}
	return strings.Join(result, ",")
}

func formatNumberToString(num uint64) string {
	units := []string{"", "Thousand", "Million", "Billion", "Trillion"}

	val := num
	result := []string{}

	i := 0
	for val > 0 {
		remainder := val % THOUSAND
		if remainder > 0 {
			result = append([]string{strings.Trim(fmt.Sprintf("%d %s", remainder, units[i]), " ")}, result...)
		}
		val = val / THOUSAND
		i += 1
	}
	return strings.Join(result, " ")
}

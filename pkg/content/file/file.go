package file

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func GetContent(desc, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "error opening file: %s", path)
	}
	defer f.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(f)

	content.WriteString(desc)
	content.WriteString("\n")

	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", errors.Wrapf(err, "error scanning file: %s", path)
	}

	return content.String(), nil
}

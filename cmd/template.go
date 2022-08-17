package cmd

import (
	"bufio"
	"strings"
)

const (
	PADDING    = "   "
	NEWLINE    = "\r\n"
	TrimPrefix = "\t\r\n"
)

func Examples(txt string) string {
	var sb strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(txt))
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			sb.WriteString(PADDING)
			sb.WriteString(strings.TrimLeft(line, TrimPrefix))
			sb.WriteString(NEWLINE)
		}
	}

	return strings.TrimRight(sb.String(), NEWLINE)
}

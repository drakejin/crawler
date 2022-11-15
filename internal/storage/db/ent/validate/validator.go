package validate

import (
	"errors"
	"github.com/drakejin/crawler/internal/_const"
	"strings"
	"unicode/utf8"
)

// MaxRuneCount validates the rune length of a string by using the unicode/utf8 package.
func MaxRuneCount(maxLen int) func(s string) error {
	return func(s string) error {
		if utf8.RuneCountInString(s) > maxLen {
			return errors.New("value is more than the max length")
		}
		return nil
	}
}

func FormatColumnsToEraseUnnecessaryCharacters(columnsStr string) string {
	columns := strings.Split(columnsStr, _const.SplitDelimiter)
	var items []string
	for _, column := range columns {
		if column != "" {
			items = append(items, strings.ReplaceAll(strings.TrimSpace(column), "\t\r", ""))
		}
	}

	return strings.Join(items, _const.SplitDelimiter)
}

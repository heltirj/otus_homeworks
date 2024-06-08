package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder
	runesInput := []rune(input)

	i := 0
	end := len(runesInput) - 1
	for i <= end {
		if unicode.IsLetter(runesInput[i]) {
			if i == end || isRuneLetterOrSlash(runesInput[i+1]) {
				result.WriteRune(runesInput[i])
				i++
				continue
			}

			if unicode.IsDigit(runesInput[i+1]) {
				s, err := repeatRune(runesInput[i], runesInput[i+1])
				if err != nil {
					return "", err
				}

				result.WriteString(s)
				i += 2
				continue
			}

			return "", ErrInvalidString
		}

		if runesInput[i] == '\\' {
			if i == end || unicode.IsLetter(runesInput[i+1]) {
				return "", ErrInvalidString
			}

			if unicode.IsDigit(runesInput[i+1]) || runesInput[i+1] == '\\' {
				if i < end-1 && unicode.IsDigit(runesInput[i+2]) {
					s, err := repeatRune(runesInput[i+1], runesInput[i+2])
					if err != nil {
						return "", err
					}

					result.WriteString(s)
					i += 3
					continue
				}

				result.WriteRune(runesInput[i+1])
				i += 2
				continue
			}
		}

		return "", ErrInvalidString
	}

	return result.String(), nil
}

func repeatRune(toRepeat, counter rune) (string, error) {
	count, err := strconv.Atoi(string(counter))
	if err != nil {
		return "", err
	}

	return strings.Repeat(string(toRepeat), count), nil
}

func isSlash(r rune) bool {
	return r == '\\'
}

func isRuneLetterOrSlash(r rune) bool {
	return unicode.IsLetter(r) || isSlash(r)
}

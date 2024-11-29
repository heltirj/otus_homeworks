package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString              = errors.New("invalid string")
	ErrWrongArgsCount             = errors.New("wrong args count")
	ErrFirstArgIsNotALetter       = errors.New("first arg is not a letter")
	ErrFirstArgIsNotADigitOrSlash = errors.New("first arg is not a digit or a slash")
)

func Unpack(input string) (string, error) {
	var result strings.Builder
	runesInput := []rune(input)

	i := 0
	end := len(runesInput) - 1
	for i <= end {
		if unicode.IsLetter(runesInput[i]) {
			if i == end {
				result.WriteRune(runesInput[i])
				break
			}

			seq, offset, err := createSequenceFromLetter(runesInput[i], runesInput[i+1])
			if err != nil {
				return "", ErrInvalidString
			}
			result.WriteString(seq)
			i += offset
			continue
		}

		if isSlash(runesInput[i]) {
			if i == end || unicode.IsLetter(runesInput[i+1]) {
				return "", ErrInvalidString
			}

			args := []rune{runesInput[i+1]}
			if i < end-1 {
				args = append(args, runesInput[i+2])
			}

			seq, offset, err := createSequenceAfterEscapeSym(args...)
			if err != nil {
				return "", ErrInvalidString
			}
			result.WriteString(seq)
			i += offset
			continue
		}

		return "", ErrInvalidString
	}

	return result.String(), nil
}

func repeatRune(toRepeat, counter rune) (string, error) {
	count, err := strconv.Atoi(string(counter))
	if err != nil {
		return "", ErrInvalidString
	}

	return strings.Repeat(string(toRepeat), count), nil
}

func isSlash(r rune) bool {
	return r == '\\'
}

func isRuneLetterOrSlash(r rune) bool {
	return unicode.IsLetter(r) || isSlash(r)
}

func createSequenceFromLetter(letter, nextRune rune) (string, int, error) {
	if !unicode.IsLetter(letter) {
		return "", 0, ErrFirstArgIsNotALetter
	}
	if isRuneLetterOrSlash(nextRune) {
		return string(letter), 1, nil
	}

	seq, err := repeatRune(letter, nextRune)
	if err != nil {
		return "", 0, err
	}
	return seq, 2, nil
}

func createSequenceAfterEscapeSym(runes ...rune) (string, int, error) {
	if len(runes) < 1 || len(runes) > 2 {
		return "", 0, ErrWrongArgsCount
	}

	if unicode.IsDigit(runes[0]) || isSlash(runes[0]) {
		if len(runes) == 2 && unicode.IsDigit(runes[1]) {
			seq, _ := repeatRune(runes[0], runes[1])
			return seq, 3, nil
		}
		return string(runes[0]), 2, nil
	}

	return "", 0, ErrFirstArgIsNotADigitOrSlash
}

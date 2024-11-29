package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	pre = regexp.MustCompile(`^\p{P}+|\p{P}+$`)
	hre = regexp.MustCompile(`^\-{2,}$`)
)

func Top10(input string) []string {
	return topKeys(getWordStats(input), 10)
}

func getWordStats(text string) map[string]int {
	if text == "" {
		return nil
	}

	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	split := strings.Split(text, " ")

	results := make(map[string]int)
	for i := range split {
		if split[i] == "" {
			continue
		}

		word := normalizeWord(split[i])

		if _, ok := results[word]; !ok {
			results[split[i]] = 1
			continue
		}
		results[word]++
	}

	return results
}

func topKeys(m map[string]int, count int) []string {
	if m == nil {
		return nil
	}

	type pair struct {
		key string
		val int
	}

	kvList := make([]pair, 0, len(m))
	for k, v := range m {
		kvList = append(kvList, pair{k, v})
	}

	sort.Slice(kvList, func(i, j int) bool {
		if kvList[i].val == kvList[j].val {
			return kvList[i].key < kvList[j].key
		}

		return kvList[i].val > kvList[j].val
	})

	if count > len(kvList) {
		count = len(kvList)
	}

	res := make([]string, count)
	for i := 0; i < count; i++ {
		res[i] = kvList[i].key
	}

	return res
}

func normalizeWord(word string) string {
	if hre.MatchString(word) {
		return word
	}
	return pre.ReplaceAllString(word, "")
}

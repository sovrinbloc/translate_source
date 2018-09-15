package conversions

import "sort"

func MapToSliceDesc(hanWords map[string]string) []string {
	lengths := make(map[int][]string)
	// each index is a slice of strings. if there are 2 of the same length, they both will be in lengths[2]{a,b}
	keys := []int{}
	// holds the length of the word
	for word, _ := range hanWords {
		keys = append(keys, len([]rune(word))) // all the lengths of each word. thats it.
		lengths[len([]rune(word))] = append(lengths[len([]rune(word))], word)
	}
	sort.Ints(keys)
	for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
		keys[i], keys[j] = keys[j], keys[i]
	}
	words := []string{}
	//this is in order
	for _, value := range keys {
		for _, word := range lengths[value] {
			words = append(words, word)
		}
	}
	return words
}

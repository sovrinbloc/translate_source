package source

import (
	"fmt"
	"regexp"
	"sort"
)

type SourceRegex struct {
	*regexp.Regexp
	source string
	words  map[string]string
	Regexs map[string]*RegexPool
}

type RegexPool struct {
	translation string
	*regexp.Regexp
	count int
}

type TranslateSource interface {
	TranslateSource(string) string
	GetStrings() string
}

func (s *SourceRegex) TranslateSource(filename string) {

}

const (
	DIR_BASE           = "/Users/joealai/go/src/github.com/shen100/golang123/"
	LOC_WEBSITE_HEADER = "website/components/Header.vue"
	LOC_MAIN           = "main.go"
)

func (j *SourceRegex) HanFind(source string) map[string]string {
	hanWords := make(map[string]string)
	allHan := j.FindAll([]byte(source), -1)
	for _, han := range allHan {
		hanWords[string(han)] = ""
		j.words[string(han)] = ""
	}

	lengths := make(map[int][]string) // each index is a slice of strings. if there are 2 of the same length, they both will be in lengths[2]{a,b}
	keys := []int{} // holds the length of the word
	for word, _ := range hanWords {
		keys = append(keys, len([]rune(word))) // all the lengths of each word. thats it.
		lengths[len([]rune(word))] = append(lengths[len([]rune(word))], word)
	}
	sort.Ints(keys)


	//this is in order
	results := make(map[string]string)
	for _, value := range keys {
		for _, word := range lengths[value] {
			results[word] = ""
		}
	}

	//but still, maps do not come in order. so we must re-order them.

	return results
}

func (j *SourceRegex) HanCreateRegexs() map[string]*RegexPool {
	j.Regexs = make(map[string]*RegexPool)
	for han, _ := range j.words {
		if k, err := regexp.Compile(string(han)); err == nil {
			if val, ok := j.Regexs[string(han)]; ok {
				val.count++
			}
			j.Regexs[string(han)] = &RegexPool{Regexp: k, count: 0}
		} else {
			panic("Error establishing something")
		}
	}
	return j.Regexs
}

func (j *SourceRegex) HanReplace(source string, replacement string) {
	r := j.ReplaceAll([]byte(source), []byte(replacement))
	fmt.Println(r)
	q := j.Find([]byte(">话题</a><"))
	q = j.ReplaceAll([]byte(">话题</a><"), []byte("go"))
	fmt.Println(string(q))
	fmt.Println(regexp.MatchString("\\p{Han}+", ">话题</a><"))
}

func NewHanRegex() *SourceRegex {
	if j, err := regexp.Compile("\\p{Han}+"); err == nil {
		return &SourceRegex{Regexp: j,
			words: make(map[string]string),
		}
	}
	return nil
}

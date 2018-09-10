package source

import (
	"errors"
	"fmt"
	"github.com/bas24/googletranslatefree"
	"testing"
)

func TestRegexFind(t *testing.T) {
	r := NewHanRegex()
	fmt.Println(r.HanFind(">话题</a><>话题</a><>话题</a><"))
	j := r.HanCreateRegexs()
	all := j["话题"].Regexp.FindAll([]byte(">话题</a><>话题</a><>话题</a><"), -1)
	translateWord(string(all[0]))
	for _, single := range all {
		translateWord(string(single))
	}

}

func translateWord(word string) {
	word, err := translategooglefree.Translate(word+"!",
		"zh",
		"en")
	if err != nil {
		errors.New(err.Error())
	}
	fmt.Println(word)
}

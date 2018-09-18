package source

import (
	"errors"
	"fmt"
	"github.com/bas24/googletranslatefree"
	"reflect"
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

func TestHanFinder(t *testing.T) {
	r := NewHanRegex()
	if got, expected := r.HanFind("Testing话题TSDF*#"), "话题"; got[0] != expected {
		t.Errorf("incorrect return value, got %s, expected %s", got, expected)
	}
	if got, expected := reflect.TypeOf(r.HanFind("Testing话题TSDF*#")), reflect.TypeOf([]string{"话题"}); got != expected {
		t.Errorf("incorrect return value, got %s, expected %s", got, expected)
	}
	if got, expected := len(r.HanFind("Testing话题TSDF*#话题sdsds")), int(2); got != expected {
		t.Errorf("incorrect return value, got %s, expected %s", got, expected)
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

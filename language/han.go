package language

import (
	"fmt"
	"github.com/dafanasev/go-yandex-translate"
	"translate_source/directory"
	"translate_source/source"
	"translate_source/config"
)

type TranslateSource struct {
	SourceFiles    *directory.LocationScan
	RegexTranslate *source.SourceRegex
	TotalSource map[string]map[string]string
	*translate.Translator

}

func (t *TranslateSource) TranslateSource(string) string {
	panic("implement me")
}

func (t *TranslateSource) GetForeignStrings() (map[string]map[string]string, map[string]*source.RegexPool) {
	dat := directory.NewLocationScan()
	dat.AddIgnoreFile("en_US", ".sql", ".key", "simplemde.js")
	dat.AddWhitelistFile(".go")
	dat.AddDirectory(source.DIR_BASE)

	h := source.NewHanRegex()

	totalSource := make(map[string]map[string]string)
	for index, file := range dat.Files {
		totalSource[index] = h.HanFind(file) // map of Han => English
	}

	regexs := h.HanCreateRegexs()

	t.SourceFiles = dat
	t.TotalSource = totalSource
	return totalSource, regexs
}

type Translate interface {
	TranslateDirectory()
	TranslateWord(string)
}

func NewTranslate() *TranslateSource {
	t := TranslateSource{Translator: translate.New(config.Env["YANDEX_KEY"])}
	return &t
}

func (tr *TranslateSource) TranslateString(source string) string {

	_, err := tr.GetLangs(config.Env["FROM"])
	if err != nil {
		fmt.Println(err)
	}

	translation, err := tr.Translate(config.Env["TO"], source)
	if err != nil {
		fmt.Println(err)
	} else {
		return translation.Result()
	}
	return ""
}

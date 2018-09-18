package language

import (
	"github.com/dafanasev/go-yandex-translate"
	"log"
	"regexp"
	"strings"
	"translate_source/config"
	"translate_source/directory"
	"translate_source/persist"
	"translate_source/source"
)

const YANDEX_MAX_LENGTH = 10000

type TranslateSource struct {
	SourceFiles    *directory.LocationScan
	RegexTranslate *source.SourceRegex
	TotalSource    map[string][]string
	*translate.Translator
}

type WordsSource struct {
	Words  []string
	Source string
}

func (t *TranslateSource) GetForeignStrings() (map[string][]string, map[string]*source.RegexPool) {
	dat := directory.NewLocationScan(true)
	if config.Env.Ignore != nil {
		dat.AddIgnoreFile(config.Env.Ignore...)
	}
	if config.Env.Whitelist != nil {
		dat.AddWhitelistFile(config.Env.Whitelist...)
	}
	dat.AddDirectory(source.DIR_BASE) // gets files and populates them with source as well
	dat.GetSources()

	h := source.NewHanRegex()

	totalSource := make(map[string][]string)
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
	t := TranslateSource{Translator: translate.New(config.Env.Vars["YANDEX_KEY"])}
	return &t
}

func (tr *TranslateSource) TranslateString(source string) string {
	_, err := tr.GetLangs(config.Env.Vars["FROM"])
	if err != nil {
		log.Println(err)
	}

	if t, err := persist.RedisClient.Get(source).Result(); err == nil {
		log.Println("translated from cache")
		return t
	}
	translation, err := tr.Translate(config.Env.Vars["TO"], source)
	if err != nil {
		log.Println(err)
	} else {
		err := persist.RedisClient.Set(source, translation.Result(), 0).Err()
		if err != nil {
			log.Printf("error caching %s as %s: %s", source, translation.Result(), err)
		}
		return translation.Result()
	}
	return ""
}

func (tr *TranslateSource) TranslateBulkString(source string) string {
	_, err := tr.GetLangs(config.Env.Vars["FROM"])
	if err != nil {
		log.Println(err)
	}

	if t, err := persist.RedisClient.Get(source).Result(); err == nil {
		log.Println("translated from cache")
		return t
	}

	sliceOfChineseWords := strings.Split(source, config.Env.Vars["DELIMITER"])
	for _, word := range sliceOfChineseWords {
		if e, err := persist.RedisClient.Get(word).Result(); err == nil {
			source = strings.Replace(source, word, e, -1)
			log.Println("translated from cache")
		}
	}

	j, err := regexp.Compile("\\p{Han}+")
	if found := j.FindIndex([]byte(source)); found == nil {
		return source
	}

	if len(source) > YANDEX_MAX_LENGTH {
		return source
	}

	translation, err := tr.Translate(config.Env.Vars["TO"], source)
	results := translation.Result()
	sliceOfEnglishWords := strings.Split(results, config.Env.Vars["DELIMITER"])
	for key, han := range sliceOfChineseWords {
		if _, err := persist.RedisClient.Get(han).Result(); err != nil {
			persist.RedisClient.Set(han, sliceOfEnglishWords[key], 0)
			log.Printf("set key for word %s: %s", han, sliceOfEnglishWords[key])
		}
	}
	if err != nil {
		log.Println(err)
	} else {
		err := persist.RedisClient.Set(source, translation.Result(), 0).Err()
		if err != nil {
			log.Printf("error caching %s as %s: %s", source, translation.Result(), err)
		}
		return translation.Result()
	}
	return ""
}

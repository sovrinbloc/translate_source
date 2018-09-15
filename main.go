package main

import (
	"time"
	"translate_source/config"
	"translate_source/directory"
	"translate_source/language"
	"strings"
	"log"
	"fmt"
)

func init() {
	language.InitRedis()
	log.Println("starting .env configuration")
	config.EnvInit()
	log.Println("successfully imported .env file")
}

func main() {
	//return
	directory.Copy(config.Env.Vars["DIRECTORY"], config.Env.Vars["NEW_DIRECTORY"])
	d := directory.NewFileStructure(config.Env.Vars["DIR_BASE"])
	d.ClipBasePath()
	d.CreateFolderStructure(config.Env.Vars["NEW_DIR"])
	translator := TranslateSourceByJoin()
	//translator := TranslateSource()
	d.AddFiles(translator.SourceFiles.Files)
	d.SaveFiles(config.Env.Vars["NEW_DIR"])

	log.Println("completed writing")
	for key, _ := range translator.SourceFiles.Files {
		log.Println(key)
	}
	log.Println("done")
}

func TranslateSource() *language.TranslateSource {
	translator := language.NewTranslate()
	totalSource, regexPool := translator.GetForeignStrings()
	for location, stringSlice := range totalSource {
		log.Println("translating location:", location)
		for index, words := range stringSlice {
			totalSource[location][index] = translator.TranslateString(words)
			translator.SourceFiles.Files[location] = string(regexPool[words].ReplaceAll(
				[]byte(translator.SourceFiles.Files[location]),
				[]byte(totalSource[location][index])))
			time.Sleep(time.Second / 100)
		}
		log.Println("finished translating", len(stringSlice), "words")
	}
	return translator
}

func TranslateSourceByJoin() *language.TranslateSource {
	translator := language.NewTranslate()
	totalSource, regexPool := translator.GetForeignStrings()

	var charCount int
	var wordCount int

	start := time.Now()

	// chineseSlice {}string{在线图书, 提示标题, 回复了你}
	for location, chineseSlice := range totalSource {
		log.Println("translating location:", location)

		if chineseWordsString := strings.Join(chineseSlice, config.Env.Vars["DELIMITER"]); len(chineseWordsString) > 0 {
			charCount += len(chineseWordsString)
			englishTranslations := translator.TranslateBulkString(chineseWordsString)
			fmt.Println(englishTranslations)

			sliceOfChineseWords := strings.Split(chineseWordsString, config.Env.Vars["DELIMITER"])
			sliceOfEnglishWords := strings.Split(englishTranslations, config.Env.Vars["DELIMITER"])

			wordCount += len(sliceOfChineseWords)

			for key, value := range sliceOfChineseWords {
				totalSource[location][key] = sliceOfEnglishWords[key]
				translator.SourceFiles.Files[location] = string(regexPool[value].ReplaceAll(
					[]byte(translator.SourceFiles.Files[location]),
					[]byte(totalSource[location][key])))
			}
			log.Printf("finished translating %v chinese words\n", len(chineseSlice))
		}
	}
	log.Printf("translated %v characters in %v words\n", charCount, wordCount)
	elapsed := time.Since(start)
	log.Printf("translation took %s at %f words/second\n", elapsed, float64(wordCount)/elapsed.Seconds())
	log.Println("finished translating successfully")
	return translator
}



package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"translate_source/config"
	"translate_source/directory"
	"translate_source/language"
	"translate_source/persist"
)

func init() {
	persist.InitRedis()
	log.Println("starting .env configuration")
	config.EnvInit()
	log.Println("successfully imported .env file")
	persist.InitRoach()
	persist.InitRedis()
	persist.CacheTranslations()
}

func main() {
	return
	directory.Copy(config.Env.Vars[config.DIRECTORY], config.Env.Vars[config.NEW_DIR])
	d := directory.NewFileStructure(config.Env.Vars[config.DIR_BASE])
	d.ClipBasePath()
	d.CreateFolderStructure(config.Env.Vars[config.NEW_DIR])
	translator := TranslateSourceByJoin()
	//translator := TranslateSource()
	d.AddFiles(translator.SourceFiles.Files)
	d.SaveFiles(config.Env.Vars[config.NEW_DIR])

	log.Println("completed writing")
	for key, _ := range translator.SourceFiles.Files {
		log.Println(key)
	}
	log.Println("done")
}

func TranslateSource() *language.TranslateSource {
	translator := language.NewTranslate()
	totalSource, regexPool := translator.GetForeignStrings()
	for location, chineseStrings := range totalSource {
		log.Println("translating location:", location)
		for index, hanWords := range chineseStrings {
			totalSource[location][index] = translator.TranslateString(hanWords)
			translator.SourceFiles.Files[location] = string(regexPool[hanWords].ReplaceAll(
				[]byte(translator.SourceFiles.Files[location]),
				[]byte(totalSource[location][index])))
			persist.TempDictionary[hanWords] = &totalSource[location][index]
			time.Sleep(time.Second / 100)
		}
		log.Println("finished translating", len(chineseStrings), "hanWords")
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
				if _, ok := persist.FullDictionary[value]; !ok {
					persist.TempDictionary[value] = &totalSource[location][key]
					persist.FullDictionary[value] = struct{}{}
				}
			}
			log.Printf("finished translating %v chinese words\n", len(chineseSlice))
			persist.SaveToDictionary()
			log.Println("persisted to dictionary")
		}
	}
	log.Printf("translated %v characters in %v words\n", charCount, wordCount)

	elapsed := time.Since(start)
	log.Printf("translation took %s at %f words/second\n", elapsed, float64(wordCount)/elapsed.Seconds())
	log.Println("finished translating successfully")
	return translator
}

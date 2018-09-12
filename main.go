package main

import (
	"fmt"
	"time"
	"translate_source/config"
	"translate_source/directory"
	"translate_source/language"
)

//todo: Add redis caching for words
//todo: Add multiple word translations at once

func init() {
	fmt.Println("Starting .env Configuration")
	config.EnvInit()
	fmt.Println("Successfully Imported .env File")
}

func main() {

	d := directory.NewFileStructure(config.Env.Vars["DIR_BASE"])
	d.ClipBasePath()
	d.CreateFolderStructure("golang123/source")

	translator := language.NewTranslate()
	totalSource, regexPool := translator.GetForeignStrings()

	for location, stringSlice := range totalSource {
		fmt.Println("Working on location:", location)
		for index, words := range stringSlice {
			totalSource[location][index] = translator.TranslateString(words)
			translator.SourceFiles.Files[location] = string(regexPool[words].ReplaceAll(
				[]byte(translator.SourceFiles.Files[location]),
				[]byte(totalSource[location][index])))
			time.Sleep(time.Second / 100)
		}
		fmt.Println("finished translating", len(stringSlice), "words")
	}
	d.AddFiles(translator.SourceFiles.Files)
	d.SaveFiles("golang123/source")

	fmt.Println("Completed writing")
	for key, _ := range translator.SourceFiles.Files {
		fmt.Println(key)
	}
	fmt.Println("Done")
}

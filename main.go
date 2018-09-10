package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
	"translate_source/config"
	"translate_source/directory"
	"translate_source/language"
	"translate_source/source"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Env["YANDEX_KEY"] = os.Getenv("YANDEX_KEY")
	config.Env["DIR_BASE"] = os.Getenv("DIRECTORY_BASE")
	config.Env["FROM"] = os.Getenv("FROM")
	config.Env["TO"] = os.Getenv("TO")
	fmt.Println("configuration finished")
}

func main() {

	translator := language.NewTranslate()
	totalSource := translator.GetForeignStrings()
	for location, value := range totalSource {
		fmt.Println(location)
		for words, _ := range value {
			//fmt.Println(location, words, value2)
			totalSource[location][words] = translator.TranslateString(words)
			fmt.Println(words, ": ", totalSource[location][words])
			//translator.SourceFiles.Files[location] = string(translator.RegexTranslate.Regexs[words].ReplaceAll(
			//	[]byte(translator.SourceFiles.Files[location]),
			//	[]byte(totalSource[location][words])))
			time.Sleep(time.Second / 10)
		}
	}

	for key, val := range translator.SourceFiles.Files {
		fmt.Println(key, val)
	}
	return
	//err := hanReplace()
	//dat, err := directory.NewLocationScan(source.DIR_BASE)
	dat := directory.NewLocationScan()
	dat.AddIgnoreFile("en_US", ".sql", ".location", "simplemde.js")
	dat.AddWhitelistFile(".go", ".vue")
	dat.AddDirectory(source.DIR_BASE)
	for _, file := range dat.GetFileList() {
		fmt.Println("Filename:", file)
	}
	for _, dir := range dat.GetDirectoryList() {
		fmt.Println("Directory:", dir)
	}
	for file, code := range dat.GetSources() {
		fmt.Println("Directory:", file)
		fmt.Println("Source:", code)
	}
}

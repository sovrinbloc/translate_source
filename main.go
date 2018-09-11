package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
	"translate_source/config"
	"translate_source/language"
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
	totalSource, regexPool := translator.GetForeignStrings()

	for location, value := range totalSource {
		for words, _ := range value {
			fmt.Println(location)
			totalSource[location][words] = translator.TranslateString(words)
			translator.SourceFiles.Files[location] = string(regexPool[words].ReplaceAll(
				[]byte(translator.SourceFiles.Files[location]),
				[]byte(totalSource[location][words])))
			fmt.Println(words, ": ", totalSource[location][words])
			fmt.Println(translator.SourceFiles.Files[location])
			time.Sleep(time.Second / 100)
		}
		fmt.Println(translator.SourceFiles.Files[location])
	}

	for key, val := range translator.SourceFiles.Files {
		fmt.Println(key, val)
	}
}

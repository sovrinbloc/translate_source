package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
	"translate_source/config"
	"translate_source/language"
	"translate_source/directory"
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

	d := directory.NewFileStructure(config.Env["DIR_BASE"])
	d.ClipBasePath()
	d.CreateFolderStructure("golang123/source")



	translator := language.NewTranslate()
	totalSource, regexPool := translator.GetForeignStrings()

	for location, stringSlice := range totalSource {
		for index, words := range stringSlice {
			fmt.Println(location)
			totalSource[location][index] = translator.TranslateString(words)
			translator.SourceFiles.Files[location] = string(regexPool[words].ReplaceAll(
				[]byte(translator.SourceFiles.Files[location]),
				[]byte(totalSource[location][index])))
			fmt.Println(words, ": ", totalSource[location][index])
			fmt.Println(translator.SourceFiles.Files[location])
			time.Sleep(time.Second / 100)
		}
		fmt.Println(translator.SourceFiles.Files[location])
	}

	d.AddFiles(translator.SourceFiles.Files)
	d.SaveFiles("golang123/source")

	for key, val := range translator.SourceFiles.Files {
		fmt.Println(key, val)
	}
}

package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"unicode"
)

var Env = struct {
	Vars      map[string]string
	Ignore    []string
	Whitelist []string
}{
	make(map[string]string),
	make([]string, 0),
	make([]string, 0),
}

//var Env.Vars = make(map[string]string)

func EnvInit() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Env.Vars["YANDEX_KEY"] = os.Getenv("YANDEX_KEY")
	Env.Vars["DIR_BASE"] = os.Getenv("DIRECTORY_BASE")
	Env.Vars["FROM"] = os.Getenv("FROM")
	Env.Vars["TO"] = os.Getenv("TO")
	Env.Vars["DELIMITER"] = os.Getenv("DELIMITER")
	Env.Vars["NEW_DIR"] = os.Getenv("NEW_DIR")

	Env.Whitelist = nil
	if whitelist := strings.Split(RemoveSpaces(os.Getenv("WHITELIST")), ","); len(whitelist) > 0 {
		Env.Whitelist = whitelist
	}
	Env.Ignore = nil
	if ignore := strings.Split(RemoveSpaces(os.Getenv("IGNORE")), ","); len(ignore) > 0 {
		Env.Ignore = ignore
	}
	fmt.Println("configuration finished")

}

func RemoveSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

package config

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
	"fmt"
	"log"
)

var Env = struct {
	Vars map[string]string
	Ignore []string
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
	Env.Whitelist = strings.Split(strings.Replace(os.Getenv("WHITELIST"), " ", "", -1), ", ")
	Env.Ignore = strings.Split(strings.Replace(os.Getenv("IGNORE"), " ", "", -1), ", ")
	fmt.Println("configuration finished")
}
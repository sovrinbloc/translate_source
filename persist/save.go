package persist

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"log"
	"regexp"
	"strings"
	"time"
	c "translate_source/config"
)

var RoachDB = new(sql.DB)
var TempDictionary = make(map[string]*string)
var FullDictionary = make(map[string]struct{})

func NewPersistDictionary() *PersistDictionary {
	redis := InitRedis()
	roach := InitRoach()
	return &PersistDictionary{
		redis,
		make(map[string]*string),
		roach,
	}
}

type PersistDictionary struct {
	*redis.Client
	Dictionary map[string]*string
	*sql.DB
}

func InitRoach() *sql.DB {
	var err error
	RoachDB, err = sql.Open("postgres",
		fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable",
			c.Env.Vars["ROACH_USER"],
			c.Env.Vars["ROACH_HOST"],
			c.Env.Vars["ROACH_PORT"],
			c.Env.Vars["ROACH_DB"]),
	)
	log.Println("CockroachDB Initialized")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	return RoachDB
}

func SaveToDictionary() {

	query := "INSERT INTO han_translations (word, translation) VALUES "
	var cleanTranslation string
	for word, translation := range TempDictionary {
		h, _ := regexp.Compile("\\w")
		if !h.Match([]byte(word)) {
			word = strings.Replace(word, "'", "", -1)
			cleanTranslation = strings.Replace(*translation, "'", "", -1)
			query = fmt.Sprintf("%s ('%s', '%s'),", query, word, cleanTranslation)
		}
	}
	query = fmt.Sprintf("%s;", query)

	query = strings.Replace(query, ",;", ";", -1)
	fmt.Println(query)
	if query != "INSERT INTO han_translations (word, translation) VALUES ;" {
		if result, err := RoachDB.Exec(query); err != nil {
			log.Fatal(err)
			fmt.Println(result.RowsAffected())
		}
	}
	time.Sleep(time.Second / 5)

	// Print out the balances.
	rows, err := RoachDB.Query("SELECT word, translation FROM han_translations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("Translated Words:")
	for rows.Next() {
		var hanWord, englishTranslation string
		if err := rows.Scan(&hanWord, &englishTranslation); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s\n", hanWord, englishTranslation)
	}
	ClearDictionary()
}

func CreateDatabase() {
	if _, err := RoachDB.Exec(
		"CREATE TABLE IF NOT EXISTS han_translations (word VARCHAR PRIMARY KEY, translation VARCHAR )"); err != nil {
		log.Fatal(err)
	}
}

func ClearDatabase() {
	// Create the "accounts" table.
	if _, err := RoachDB.Exec(
		"DROP TABLE IF EXISTS han_translations;"); err != nil {
		log.Fatal(err)
	}
}

func ClearDictionary() {
	TempDictionary = make(map[string]*string)
}

func PrintTranslations() {
	// Print out the balances.
	rows, err := RoachDB.Query("SELECT word, translation FROM han_translations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("Translated Words:")
	for rows.Next() {
		var hanWord, englishTranslation string
		if err := rows.Scan(&hanWord, &englishTranslation); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s\n", hanWord, englishTranslation)
	}
}

func CacheTranslations() {
	rows, err := RoachDB.Query("SELECT word, translation FROM han_translations")
	h, _ := regexp.Compile("\\w")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var hanWord, englishTranslation string
		if err := rows.Scan(&hanWord, &englishTranslation); err != nil {
			log.Fatal(err)
		}

		if !h.Match([]byte(hanWord)) {
			RedisClient.Set(hanWord, englishTranslation, 0)
			FullDictionary[hanWord] = struct{}{}
		}
	}
	log.Println("translations caches")
}

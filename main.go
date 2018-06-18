package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/toferc/ore_web_roller/database"
)

var db *pg.DB

func init() {
	os.Setenv("DBUser", "chris")
	os.Setenv("DBPass", "12345")
	os.Setenv("DBName", "ore_engine")

}

func main() {

	if os.Getenv("ENVIRONMENT") == "production" {
		url, ok := os.LookupEnv("DATABASE_URL")

		if !ok {
			log.Fatalln("$DATABASE_URL is required")
		}

		options, err := pg.ParseURL(url)

		if err != nil {
			log.Fatalf("Connection error: %s", err.Error())
		}

		db = pg.Connect(options)
	} else {
		db = pg.Connect(&pg.Options{
			User:     os.Getenv("DBUser"),
			Password: os.Getenv("DBPass"),
			Database: os.Getenv("DBName"),
		})
	}

	defer db.Close()

	err := database.InitDB(db)
	if err != nil {
		panic(err)
	}

	fmt.Println(db)

	// Check for terminal flag
	t := flag.Bool("t", false, "Activate the local terminal")

	flag.Parse()

	if *t == true {
		Terminal(db)
	} else {

		port := os.Getenv("PORT")

		if port == "" {
			port = "8080"
		}

		fmt.Println("Starting Webserver at port " + port)
		http.HandleFunc("/", IndexHandler)
		http.HandleFunc("/roll/", RollHandler)
		http.HandleFunc("/opposed/", OpposeHandler)
		http.HandleFunc("/view/", CharacterHandler)
		http.HandleFunc("/new/", NewCharacterHandler)
		http.HandleFunc("/modify/", ModifyCharacterHandler)
		http.HandleFunc("/delete/", DeleteCharacterHandler)

		log.Fatal(http.ListenAndServe(":"+port, nil))
	}

}

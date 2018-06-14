package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/toferc/ore_web/database"
)

var db *pg.DB

func init() {
	os.Setenv("DBUser", "chris")
	os.Setenv("DBPass", "12345")
	os.Setenv("DBName", "ore_engine")

}

func main() {

	db = pg.Connect(&pg.Options{
		User:     os.Getenv("DBUser"),
		Password: os.Getenv("DBPass"),
		Database: os.Getenv("DBName"),
	})

	defer db.Close()

	err := database.InitDB(db)
	if err != nil {
		panic(err)
	}

	fmt.Println(db)

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
		http.HandleFunc("/character/", CharacterHandler)

		log.Fatal(http.ListenAndServe(":"+port, nil))
	}

}

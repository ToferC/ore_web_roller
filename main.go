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
		http.HandleFunc("/", CharacterIndexHandler)
		http.HandleFunc("/roll/", RollHandler)
		http.HandleFunc("/opposed/", OpposeHandler)
		http.HandleFunc("/view_character/", CharacterHandler)
		http.HandleFunc("/new/", NewCharacterHandler)
		http.HandleFunc("/modify/", ModifyCharacterHandler)
		http.HandleFunc("/delete/", DeleteCharacterHandler)

		http.HandleFunc("/index_powers/", PowerIndexHandler)
		http.HandleFunc("/view_power/", PowerHandler)

		http.HandleFunc("/add_power/", AddPowerHandler)
		http.HandleFunc("/modify_power/", ModifyPowerHandler)
		http.HandleFunc("/delete_power/", DeletePowerHandler)

		http.HandleFunc("/add_hyperstat/", AddHyperStatHandler)
		http.HandleFunc("/modify_hyperstat/", ModifyHyperStatHandler)
		http.HandleFunc("/delete_hyperstat/", DeleteHyperStatHandler)

		http.HandleFunc("/add_hyperskill/", AddHyperSkillHandler)
		http.HandleFunc("/modify_hyperskill/", ModifyHyperSkillHandler)
		http.HandleFunc("/delete_hyperskill/", DeleteHyperSkillHandler)

		http.HandleFunc("/add_skill/", AddSkillHandler)
		http.HandleFunc("/add_advantages/", ModifyAdvantageHandler)

		log.Fatal(http.ListenAndServe(":"+port, nil))
	}

}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

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

		r := mux.NewRouter()

		fmt.Println("Starting Webserver at port " + port)
		r.HandleFunc("/", CharacterIndexHandler)
		r.HandleFunc("/roll/", RollHandler)
		r.HandleFunc("/opposed/", OpposeHandler)
		r.HandleFunc("/view_character/{id}", CharacterHandler)
		r.HandleFunc("/new/{setting}", NewCharacterHandler)
		r.HandleFunc("/modify/{id}", ModifyCharacterHandler)
		r.HandleFunc("/delete/{id}", DeleteCharacterHandler)

		r.HandleFunc("/index_powers/", PowerIndexHandler)
		r.HandleFunc("/view_power/{id}", PowerHandler)

		r.HandleFunc("/add_power/{id}", AddPowerHandler)
		r.HandleFunc("/add_power_from_list/{id}", PowerListHandler)
		r.HandleFunc("/modify_power/{id}/{power}", ModifyPowerHandler)
		r.HandleFunc("/delete_power/{id}/{power}", DeletePowerHandler)

		r.HandleFunc("/add_hyperstat/{id}/{stat}", AddHyperStatHandler)
		r.HandleFunc("/modify_hyperstat/{id}/{stat}", ModifyHyperStatHandler)
		r.HandleFunc("/delete_hyperstat/{id}/{stat}", DeleteHyperStatHandler)

		r.HandleFunc("/add_hyperskill/{id}/{skill}", AddHyperSkillHandler)
		r.HandleFunc("/modify_hyperskill/{id}/{skill}", ModifyHyperSkillHandler)
		r.HandleFunc("/delete_hyperskill/{id}/{skill}", DeleteHyperSkillHandler)

		r.HandleFunc("/add_skill/{id}/{skill}", AddSkillHandler)
		r.HandleFunc("/add_advantages/{id}", ModifyAdvantageHandler)

		http.Handle("/", r)

		log.Fatal(http.ListenAndServe(":"+port, r))
	}

}

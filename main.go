package main

/*

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
)

func init() {
	os.Setenv("DBUser", "chris")
	os.Setenv("DBPass", "12345")
	os.Setenv("DBName", "ore_engine")
}

func main() {

	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DBUser"),
		Password: os.Getenv("DBPass"),
		Database: os.Getenv("DBName"),
	})

	defer db.Close()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting Webserver at port " + port)
	http.HandleFunc("/roll/", RollHandler)
	http.HandleFunc("/opposed/", OpposeHandler)
	http.HandleFunc("/character/", CharacterHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
*/

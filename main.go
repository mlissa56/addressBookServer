package main

import (
	"addressBookServer/controllers/stdhttp"
	"addressBookServer/gate/psg"
	"log"
	"net/http"
	"os"

    "github.com/joho/godotenv"
)

func main() {
    enverr := godotenv.Load()
    if enverr != nil {
        log.Fatal(enverr)
    }

    pgurl := os.Getenv("DB_URL")

    db := psg.NewPsg(pgurl) 
    controller := stdhttp.NewController(db)  

    mux := http.NewServeMux()
    mux.Handle("/address-book/", controller)
    err := http.ListenAndServe(":" + os.Getenv("ADDRESSBOOK_PORT"), mux)
    if err != nil {
        log.Fatal(err)
    }            
}

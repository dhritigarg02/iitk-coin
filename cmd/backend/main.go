package main

import (
	"log"
	"net/http"
	"os"
	"github.com/dhritigarg02/iitk-coin/pkg/db"
	"github.com/dhritigarg02/iitk-coin/pkg/api"
	"github.com/gorilla/handlers"
)

func main() {

	server := api.Server{}

	server.DBstore = db.NewStore(db.ConnectDB())

	server.SetupRouter()

	log.Println("[MAIN] [INFO] Starting server at port 8080....")
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, server.Router)))
}
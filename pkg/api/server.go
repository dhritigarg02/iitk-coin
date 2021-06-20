package api

import (
	"net/http"

	"github.com/dhritigarg02/iitk-coin/pkg/db"
)

type Server struct {

	DBstore db.DBStore
	Router *http.ServeMux
}

func (server *Server) SetupRouter() {

	router := http.NewServeMux()

	router.HandleFunc("/", server.HelloHandler)
	router.HandleFunc("/login", server.Login)
	router.HandleFunc("/signup", server.Signup)
	router.HandleFunc("/secretpage", server.Secretpage)
	router.HandleFunc("/reward", server.RewardCoins)
	router.HandleFunc("/transfer", server.TransferCoins)
	router.HandleFunc("/getbalance", server.GetBalance)

	server.Router = router
}
package api

import (
	"net/http"

	"github.com/dhritigarg02/iitk-coin/pkg/db"
	"github.com/dhritigarg02/iitk-coin/pkg/middleware"
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
	router.Handle("/reward", middleware.NewEnsureAuth(server.RewardCoins))
	router.Handle("/transfer", middleware.NewEnsureAuth(server.TransferCoins))
	router.Handle("/getbalance", middleware.NewEnsureAuth(server.GetBalance))

	server.Router = router
}
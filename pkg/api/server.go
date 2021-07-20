// Package api sets up the server and contains all the handler functions.
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

func NewEnsureAuth(handlerToWrap middleware.AuthenticatedHandler) *middleware.EnsureAuth {
    return &middleware.EnsureAuth{
		Handler : handlerToWrap,
	}
}

func (server *Server) NewEnsureAdmin(handlerToWrap http.HandlerFunc) *middleware.EnsureAdmin {
    return &middleware.EnsureAdmin{
		Handler : handlerToWrap, 
		DBstore : server.DBstore,
	}
}

func (server *Server) SetupRouter() {

	router := http.NewServeMux()

	router.HandleFunc("/", server.HelloHandler)
	router.HandleFunc("/login", server.Login)
	router.HandleFunc("/signup", server.Signup)
	router.Handle("/reward", server.NewEnsureAdmin(server.RewardCoins))
	router.Handle("/transfer", NewEnsureAuth(server.TransferCoins))
	router.Handle("/getbalance", NewEnsureAuth(server.GetBalance))
	router.Handle("/redeem", NewEnsureAuth(server.RedeemCoins))

	server.Router = router
}